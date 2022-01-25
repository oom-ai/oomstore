package sqlutil

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	InsertBatchSize = 20
)

const JOIN_QUERY = `
INSERT INTO {{ qt .TableName }} ( {{ qt .EntityKey }}, {{ qt .UnixMilli }}, {{ columnJoin .Columns }})
SELECT
	l.{{ qt .EntityKey }},
	l.{{ qt .UnixMilli }},
	{{ columnJoin .Columns }}
FROM
	{{ qt .EntityRowsTableName }} AS l
LEFT JOIN {{ qt .SnapshotTable }} AS r
ON l.{{ qt .EntityKey }} = r.{{ qt .EntityName }}
WHERE l.{{ qt .UnixMilli }} >= ? AND l.{{ qt .UnixMilli }} < ?
`

const CDC_JOIN_QUERY = `
INSERT INTO {{ qt .TableName }} ( {{ qt .EntityKey }}, {{ qt .UnixMilli }}, {{ columnJoinWithComma .ValueNames }} {{ columnJoin .FeatureNames }})
SELECT
	l.{{ qt .EntityKey }},
	l.{{ qt .UnixMilli }},
	{{ columnJoinWithComma .ValueNames }}
	{{ featureValue .FeatureNames }}
FROM
	{{ qt .SnapshotJoinedTable }} AS l
LEFT JOIN {{ qt .CdcTable }} AS r
ON l.{{ qt .EntityKey }} = r.{{ qt .EntityName }} AND l.{{ qt .UnixMilli }} >= r.{{ qt .UnixMilli }}
WHERE
	l.{{ qt .UnixMilli }} >= ? AND l.{{ qt .UnixMilli }} < ? AND
	(
		r.{{ qt .UnixMilli }} IS NULL OR
		r.{{ qt .UnixMilli }} = (
			SELECT MAX({{ qt .UnixMilli }})
			FROM {{ qt .CdcTable }} AS r2
			WHERE
				l.{{ qt .EntityKey }} = r2.{{ qt .EntityName }} AND
				l.{{ qt .UnixMilli }} >= r2.{{ qt .UnixMilli }}
		)
	)
`

const READ_JOIN_RESULT_QUERY = `
SELECT
	{{ qt .EntityRowsTableName }}.{{ qt .EntityKey }},
	{{ qt .EntityRowsTableName }}.{{ qt .UnixMilli }},
	{{ fieldJoin .Fields }}
FROM {{ qt .EntityRowsTableName }}
{{ range $pair := .JoinTables }}
	{{- $t1 := qt $pair.LeftTable -}}
	{{- $t2 := qt $pair.RightTable -}}
LEFT JOIN {{ $t2 }}
ON {{ $t1 }}.{{ qt $.UnixMilli }} = {{ $t2 }}.{{ qt $.UnixMilli }} AND {{ $t1 }}.{{ qt $.EntityKey }} = {{ $t2 }}.{{ qt $.EntityKey }}
{{end}}`

func prepareJoinedTable(
	ctx context.Context,
	dbOpt dbutil.DBOpt,
	features types.FeatureList,
	groupName string,
	valueNames []string,
) (string, error) {
	// Step 1: create table joined_
	tableName := dbutil.TempTable(fmt.Sprintf("joined_%s", groupName))
	qtTableName, columnDefs, err := prepareTableSchema(dbOpt, prepareTableSchemaParams{
		tableName:    tableName,
		entityName:   "entity_key",
		valueNames:   valueNames,
		hasUnixMilli: true,
	})
	if err != nil {
		return "", err
	}
	columnFormat, err := dbutil.GetColumnFormat(dbOpt.Backend)
	if err != nil {
		return "", err
	}
	for _, f := range features {
		dbValueType, err := dbutil.DBValueType(dbOpt.Backend, f.ValueType)
		if err != nil {
			return "", err
		}
		columnDefs = append(columnDefs, fmt.Sprintf(columnFormat, f.Name, dbValueType))

	}
	schema := `
		CREATE TABLE %s (
			%s
		);
	`
	schema = fmt.Sprintf(schema, qtTableName, strings.Join(columnDefs, ",\n"))
	if err = dbOpt.ExecContext(ctx, schema, nil); err != nil {
		return "", err
	}

	qt := dbutil.QuoteFn(dbOpt.Backend)
	// Step 2: create index on table joined_
	if supportIndex(dbOpt.Backend) {
		index := fmt.Sprintf(`CREATE INDEX %s ON %s (unix_milli, entity_key)`, qt("idx_"+tableName), qtTableName)
		if err = dbOpt.ExecContext(ctx, index, nil); err != nil {
			return "", err
		}
	}
	return tableName, nil
}

func prepareEntityRowsTable(ctx context.Context,
	dbOpt dbutil.DBOpt,
	entityRows <-chan types.EntityRow,
	valueNames []string,
) (string, error) {
	// Step 1: create table entity_rows
	tableName := dbutil.TempTable("entity_rows")
	qtTableName, columnDefs, err := prepareTableSchema(dbOpt, prepareTableSchemaParams{
		tableName:    tableName,
		entityName:   "entity_key",
		valueNames:   valueNames,
		hasUnixMilli: true,
	})
	if err != nil {
		return "", err
	}
	schema := fmt.Sprintf(`
		CREATE TABLE %s (
			%s
		);
	`, qtTableName, strings.Join(columnDefs, ",\n"))

	if err = dbOpt.ExecContext(ctx, schema, nil); err != nil {
		return "", err
	}

	// Step 2: populate dataset to the table
	if err = insertEntityRows(ctx, dbOpt, tableName, entityRows, valueNames); err != nil {
		return "", err
	}

	// Step 3: create index on table entity_rows
	if supportIndex(dbOpt.Backend) {
		index := fmt.Sprintf(`CREATE INDEX idx_%s ON %s (unix_milli, entity_key)`, tableName, tableName)
		if err = dbOpt.ExecContext(ctx, index, nil); err != nil {
			return "", err
		}
	}

	return tableName, nil
}

type prepareTableSchemaParams struct {
	tableName    string
	entityName   string
	valueNames   []string
	hasUnixMilli bool
}

func prepareTableSchema(dbOpt dbutil.DBOpt, params prepareTableSchemaParams) (string, []string, error) {
	columnFormat, err := dbutil.GetColumnFormat(dbOpt.Backend)
	if err != nil {
		return "", nil, err
	}
	qt := dbutil.QuoteFn(dbOpt.Backend)

	// TODO: infer db_type from value_type
	var entityType, valueType, qtTableName string
	switch dbOpt.Backend {
	case types.BackendBigQuery:
		entityType = "STRING"
		valueType = "STRING"
		qtTableName = fmt.Sprintf("%s.%s", *dbOpt.DatasetID, qt(params.tableName))
	case types.BackendMySQL:
		entityType = "VARCHAR(255)"
		valueType = "TEXT"
		qtTableName = qt(params.tableName)
	case types.BackendSnowflake:
		entityType = "TEXT"
		valueType = "TEXT"
		qtTableName = fmt.Sprintf("PUBLIC.%s", qt(params.tableName))
	default:
		entityType = "TEXT"
		valueType = "TEXT"
		qtTableName = qt(params.tableName)
	}

	columnDefs := []string{
		fmt.Sprintf(`%s %s NOT NULL`, qt(params.entityName), entityType),
	}
	if params.hasUnixMilli {
		columnDefs = append(columnDefs, fmt.Sprintf(`%s BIGINT NOT NULL`, qt("unix_milli")))
	}
	for _, name := range params.valueNames {
		columnDefs = append(columnDefs, fmt.Sprintf(columnFormat, name, valueType))
	}

	return qtTableName, columnDefs, nil
}

func insertEntityRows(ctx context.Context,
	dbOpt dbutil.DBOpt,
	tableName string,
	entityRows <-chan types.EntityRow,
	valueNames []string,
) error {
	records := make([]interface{}, 0, InsertBatchSize)
	columns := []string{"entity_key", "unix_milli"}
	columns = append(columns, valueNames...)

	format := `%s`
	if dbOpt.Backend == types.BackendBigQuery {
		format = `"%s"`
	}
	for entityRow := range entityRows {
		record := []interface{}{fmt.Sprintf(format, entityRow.EntityKey), entityRow.UnixMilli}
		for _, v := range entityRow.Values {
			record = append(record, fmt.Sprintf(format, v))
		}
		records = append(records, record)
		if len(records) == InsertBatchSize {

			if err := dbutil.InsertRecordsToTable(ctx, dbOpt, tableName, records, columns); err != nil {
				return err
			}
			records = make([]interface{}, 0, InsertBatchSize)
		}
	}
	if err := dbutil.InsertRecordsToTable(ctx, dbOpt, tableName, records, columns); err != nil {
		return err
	}
	return nil
}

func dropTemporaryTables(ctx context.Context, db *sqlx.DB, tableNames []string) error {
	var err error
	for _, tableName := range tableNames {
		if tmpErr := dropTable(ctx, db, tableName); tmpErr != nil {
			err = tmpErr
		}
	}
	return err
}

func dropTable(ctx context.Context, db *sqlx.DB, tableName string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName)
	_, err := db.ExecContext(ctx, query)
	return errdefs.WithStack(err)
}

func supportIndex(backendType types.BackendType) bool {
	for _, b := range []types.BackendType{types.BackendSnowflake, types.BackendRedshift, types.BackendBigQuery} {
		if b == backendType {
			return false
		}
	}
	return true
}

type joinTablePair struct {
	LeftTable  string
	RightTable string
}

type readJoinResultQueryParams struct {
	EntityRowsTableName string
	EntityKey           string
	UnixMilli           string
	Fields              []string
	JoinTables          []joinTablePair
	Backend             types.BackendType
	DatasetID           string
}

func buildReadJoinResultQuery(query string, params readJoinResultQueryParams) (string, error) {
	qt := dbutil.QuoteFn(params.Backend)
	t := template.Must(template.New("temp_join").Funcs(template.FuncMap{
		"qt": qt,
		"fieldJoin": func(fields []string) string {
			return strings.Join(fields, ",\n\t")
		},
	}).Parse(query))

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, params); err != nil {
		return "", errdefs.WithStack(err)
	}
	return buf.String(), nil
}

type joinQueryParams struct {
	TableName           string
	EntityName          string
	EntityKey           string
	UnixMilli           string
	Columns             []string
	EntityRowsTableName string
	SnapshotTable       string
	Backend             types.BackendType
	DatasetID           *string
}

func buildJoinQuery(params joinQueryParams) (string, error) {
	if params.Backend == types.BackendBigQuery {
		params.TableName = fmt.Sprintf("%s.%s", *params.DatasetID, params.TableName)
		params.EntityRowsTableName = fmt.Sprintf("%s.%s", *params.DatasetID, params.EntityRowsTableName)
		params.SnapshotTable = fmt.Sprintf("%s.%s", *params.DatasetID, params.SnapshotTable)
	}

	qt := dbutil.QuoteFn(params.Backend)
	t := template.Must(template.New("join").Funcs(template.FuncMap{
		"qt": qt,
		"columnJoin": func(columns []string) string {
			return qt(columns...)
		},
	}).Parse(JOIN_QUERY))

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, params); err != nil {
		return "", errdefs.WithStack(err)
	}
	return buf.String(), nil
}

type cdcJoinQueryParams struct {
	TableName           string
	EntityName          string
	EntityKey           string
	UnixMilli           string
	ValueNames          []string
	FeatureNames        []string
	SnapshotJoinedTable string
	CdcTable            string
	Backend             types.BackendType
	DatasetID           *string
}

func buildCdcJoinQuery(params cdcJoinQueryParams) (string, error) {
	if params.Backend == types.BackendBigQuery {
		params.TableName = fmt.Sprintf("%s.%s", *params.DatasetID, params.TableName)
		params.SnapshotJoinedTable = fmt.Sprintf("%s.%s", *params.DatasetID, params.SnapshotJoinedTable)
		params.CdcTable = fmt.Sprintf("%s.%s", *params.DatasetID, params.CdcTable)
	}

	qt := dbutil.QuoteFn(params.Backend)
	t := template.Must(template.New("cdc_join").Funcs(template.FuncMap{
		"qt": qt,
		"columnJoin": func(columns []string) string {
			return qt(columns...)
		},
		"columnJoinWithComma": func(columns []string) string {
			if len(columns) == 0 {
				return ""
			}
			return fmt.Sprintf("%s,", qt(columns...))
		},
		"featureValue": func(features []string) string {
			values := make([]string, 0, len(features))
			for _, c := range features {
				values = append(values, fmt.Sprintf("(CASE WHEN r.%s IS NULL THEN l.%s ELSE r.%s END) AS %s", qt("unix_milli"), qt(c), qt(c), c))
			}
			return strings.Join(values, ",")
		},
	}).Parse(CDC_JOIN_QUERY))

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, params); err != nil {
		return "", errdefs.WithStack(err)
	}
	return buf.String(), nil
}
