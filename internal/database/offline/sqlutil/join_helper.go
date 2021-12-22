package sqlutil

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	InsertBatchSize = 20
)

func supportIndex(backendType types.BackendType) bool {
	for _, b := range []types.BackendType{types.SNOWFLAKE, types.REDSHIFT, types.BIGQUERY} {
		if b == backendType {
			return false
		}
	}
	return true
}

func createTableJoined(
	ctx context.Context,
	db *sqlx.DB,
	features types.FeatureList,
	entity types.Entity,
	groupName string,
	valueNames []string,
	backendType types.BackendType,
) (string, error) {
	columnFormat, err := dbutil.GetColumnFormat(backendType)
	if err != nil {
		return "", err
	}
	qt, err := dbutil.QuoteFn(backendType)
	if err != nil {
		return "", err
	}
	// create table joined
	tableName := dbutil.TempTable(fmt.Sprintf("joined_%s", groupName))
	columnDefs := []string{
		fmt.Sprintf(`%s  VARCHAR(%d) NOT NULL`, qt("entity_key"), entity.Length),
		fmt.Sprintf(`%s  BIGINT NOT NULL`, qt("unix_milli")),
	}
	for _, name := range valueNames {
		columnDefs = append(columnDefs, fmt.Sprintf(columnFormat, name, "TEXT"))
	}
	for _, f := range features {
		dbValueType, err := dbutil.DBValueType(backendType, f.ValueType)
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

	schema = fmt.Sprintf(schema, qt(tableName), strings.Join(columnDefs, ",\n"))
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return "", err
	}

	// snowflake doesn't support index
	if backendType != types.SNOWFLAKE {
		index := fmt.Sprintf(`CREATE INDEX idx_%s ON %s (unix_milli, entity_key)`, tableName, tableName)
		if _, err = db.ExecContext(ctx, index); err != nil {
			return "", err
		}
	}
	return tableName, nil
}

func PrepareEntityRowsTable(ctx context.Context,
	dbOpt dbutil.DBOpt,
	entity types.Entity,
	entityRows <-chan types.EntityRow,
	valueNames []string,
) (string, error) {
	// Step 1: create entity_rows table
	columnFormat, err := dbutil.GetColumnFormat(dbOpt.Backend)
	if err != nil {
		return "", err
	}
	qt, err := dbutil.QuoteFn(dbOpt.Backend)
	if err != nil {
		return "", err
	}
	tableName := dbutil.TempTable("entity_rows")

	// TODO: infer db_type from value_type
	var entityType, valueType, qtTableName string
	switch dbOpt.Backend {
	case types.BIGQUERY:
		entityType = "STRING"
		valueType = "STRING"
		qtTableName = fmt.Sprintf("%s.%s", *dbOpt.DatasetID, qt(tableName))
	default:
		entityType = fmt.Sprintf("VARCHAR(%d)", entity.Length)
		valueType = "TEXT"
		qtTableName = qt(tableName)
	}

	columnDefs := []string{
		fmt.Sprintf(`%s %s NOT NULL`, qt("entity_key"), entityType),
		fmt.Sprintf(`%s BIGINT NOT NULL`, qt("unix_milli")),
	}
	for _, name := range valueNames {
		columnDefs = append(columnDefs, fmt.Sprintf(columnFormat, name, valueType))
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
	if dbOpt.Backend == types.BIGQUERY {
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

func dropTable(ctx context.Context, db *sqlx.DB, tableName string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName)
	_, err := db.ExecContext(ctx, query)
	return err
}

const READ_JOIN_RESULT_QUERY = `
SELECT
	{{ qt .EntityRowsTableName }}.{{ .EntityKeyStr }},
	{{ qt .EntityRowsTableName }}.{{ .UnixMilliStr }},
	{{ fieldJoin .Fields }}
FROM {{ qt .EntityRowsTableName }}
{{ range $pair := .JoinTables }}
	{{- $t1 := qt $pair.LeftTable -}}
	{{- $t2 := qt $pair.RightTable -}}
lEFT JOIN {{ $t2 }}
ON {{ $t1 }}.{{ $.UnixMilliStr }} = {{ $t2 }}.{{ $.UnixMilliStr }} AND {{ $t1 }}.{{ $.EntityKeyStr }} = {{ $t2 }}.{{ $.EntityKeyStr }}
{{end}}`

type joinTablePair struct {
	LeftTable  string
	RightTable string
}

type readJoinResultQuery struct {
	EntityRowsTableName string
	EntityKeyStr        string
	UnixMilliStr        string
	Fields              []string
	JoinTables          []joinTablePair
	Backend             types.BackendType
}

func buildReadJoinResultQuery(schema readJoinResultQuery) (string, error) {
	qt, err := dbutil.QuoteFn(schema.Backend)
	if err != nil {
		return "", err
	}
	t := template.Must(template.New("temp_join").Funcs(template.FuncMap{
		"qt": qt,
		"fieldJoin": func(fields []string) string {
			return strings.Join(fields, ",\n\t")
		},
	}).Parse(READ_JOIN_RESULT_QUERY))

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, schema); err != nil {
		return "", err
	}
	return buf.String(), nil
}

const JOIN_QUERY = `
INSERT INTO {{ qt .TableName }} ( {{ .EntityKeyStr }}, {{ .UnixMilliStr }}, {{ columnJoin .Columns }})
SELECT
	l.{{ .EntityKeyStr }} AS entity_key,
	l.{{ .UnixMilliStr }} AS unix_milli,
	{{ columnJoin .Columns }}
FROM
	{{ qt .EntityRowsTableName }} AS l
LEFT JOIN {{ qt .DataTable }} AS r
ON l.{{ .EntityKeyStr }} = r.{{ qt .EntityName }}
WHERE l.{{ .UnixMilliStr }} >= ? AND l.{{ .UnixMilliStr }} < ?
`

type joinQuery struct {
	TableName           string
	EntityKeyStr        string
	EntityName          string
	UnixMilliStr        string
	Columns             []string
	EntityRowsTableName string
	DataTable           string
	Backend             types.BackendType
}

func buildJoinQuery(schema joinQuery) (string, error) {
	qt, err := dbutil.QuoteFn(schema.Backend)
	if err != nil {
		return "", err
	}

	t := template.Must(template.New("join").Funcs(template.FuncMap{
		"qt": qt,
		"columnJoin": func(columns []string) string {
			return qt(columns...)
		},
	}).Parse(JOIN_QUERY))

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, schema); err != nil {
		return "", err
	}
	return buf.String(), nil
}
