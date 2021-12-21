package bigquery

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
)

const (
	InsertBatchSize = 20
)

func createTableJoined(ctx context.Context, db *DB, features types.FeatureList, groupName string, valueNames []string, backendType types.BackendType) (string, error) {
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
		fmt.Sprintf(`%s STRING NOT NULL`, qt("entity_key")),
		fmt.Sprintf(`%s BIGINT NOT NULL`, qt("unix_milli")),
	}
	for _, name := range valueNames {
		columnDefs = append(columnDefs, fmt.Sprintf(columnFormat, name, "STRING"))
	}
	for _, f := range features {
		sqlType, err := convertValueTypeToBigQuerySQLType(f.ValueType)
		if err != nil {
			return "", err
		}
		columnDefs = append(columnDefs, fmt.Sprintf(columnFormat, f.Name, sqlType))
	}
	schema := fmt.Sprintf(`
		CREATE TABLE %s.%s (
			%s
		);
	`, db.datasetID, qt(tableName), strings.Join(columnDefs, ",\n"))
	if _, err := db.Query(schema).Read(ctx); err != nil {
		return "", err
	}

	return tableName, nil
}

func createAndImportTableEntityRows(ctx context.Context, db *DB, entityRows <-chan types.EntityRow, valueNames []string, backendType types.BackendType) (string, error) {
	columnFormat, err := dbutil.GetColumnFormat(backendType)
	if err != nil {
		return "", err
	}
	qt, err := dbutil.QuoteFn(backendType)
	if err != nil {
		return "", err
	}

	// create table entity_rows
	tableName := dbutil.TempTable("entity_rows")
	columnDefs := []string{
		fmt.Sprintf(`%s STRING NOT NULL`, qt("entity_key")),
		fmt.Sprintf(`%s BIGINT NOT NULL`, qt("unix_milli")),
	}
	for _, name := range valueNames {
		columnDefs = append(columnDefs, fmt.Sprintf(columnFormat, name, "STRING"))
	}
	schema := fmt.Sprintf(`
		CREATE TABLE %s.%s (
			%s
		);
	`, db.datasetID, qt(tableName), strings.Join(columnDefs, ",\n"))

	if _, err = db.Query(schema).Read(ctx); err != nil {
		return "", err
	}

	// populate dataset to the table
	if err := insertEntityRows(ctx, db, tableName, entityRows, valueNames, backendType); err != nil {
		return "", err
	}

	return tableName, nil
}

func insertEntityRows(ctx context.Context, db *DB, tableName string, entityRows <-chan types.EntityRow, valueNames []string, backendType types.BackendType) error {
	records := make([][]interface{}, 0, InsertBatchSize)
	columns := []string{"entity_key", "unix_milli"}
	columns = append(columns, valueNames...)
	for entityRow := range entityRows {
		record := []interface{}{fmt.Sprintf(`"%s"`, entityRow.EntityKey), entityRow.UnixMilli}
		for _, v := range entityRow.Values {
			record = append(record, fmt.Sprintf(`"%s"`, v))
		}
		records = append(records, record)
		if len(records) == InsertBatchSize {
			if err := insertRecordsToTable(db, ctx, tableName, records, columns, backendType); err != nil {
				return err
			}
			records = make([][]interface{}, 0, InsertBatchSize)
		}
	}
	if err := insertRecordsToTable(db, ctx, tableName, records, columns, backendType); err != nil {
		return err
	}
	return nil
}

func dropTable(ctx context.Context, db *DB, tableName string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName)
	_, err := db.Query(query).Read(ctx)
	return err
}

func insertRecordsToTable(db *DB, ctx context.Context, tableName string, records [][]interface{}, columns []string, backendType types.BackendType) error {
	query, err := buildQueryForInsertRecords(db.datasetID, tableName, records, columns, backendType)
	if err != nil {
		return err
	}
	if query == "" {
		return nil
	}
	if _, err := db.Query(query).Read(ctx); err != nil {
		return err
	}
	return nil
}

func buildQueryForInsertRecords(datasetID, tableName string, records [][]interface{}, columns []string, backendType types.BackendType) (string, error) {
	if len(records) == 0 {
		return "", nil
	}
	values := make([]string, 0, len(records))
	for _, row := range records {
		values = append(values, fmt.Sprintf("(%s)", strings.Join(cast.ToStringSlice(row), ",")))
	}
	columnStr := dbutil.Quote("`", columns...)
	tableName = fmt.Sprintf("`%s`", tableName)

	return fmt.Sprintf(`INSERT INTO %s.%s (%s) VALUES %s`, datasetID, tableName, columnStr, strings.Join(values, ",")), nil
}

func convertValueTypeToBigQuerySQLType(t string) (string, error) {
	switch t {
	case types.STRING:
		return "STRING", nil
	case types.INT64:
		return "BIGINT", nil
	case types.BOOL:
		return "BOOL", nil
	case types.FLOAT64:
		return "FLOAT64", nil
	case types.BYTES:
		return "BYTES", nil
	case types.TIME:
		return "DATETIME", nil
	default:
		return "", fmt.Errorf("unsupported value type %s", t)
	}
}

const READ_JOIN_RESULT_QUERY = `
SELECT
	{{ qt .EntityRowsTableName }}.{{ .EntityKeyStr }},
	{{ qt .EntityRowsTableName }}.{{ .UnixMilliStr }},
	{{ fieldJoin .Fields }}
FROM {{ $.DatasetID }}.{{ qt .EntityRowsTableName }}
{{ range $pair := .JoinTables }}
	{{- $t1 := qt $pair.LeftTable -}}
	{{- $t2 := qt $pair.RightTable -}}
lEFT JOIN {{ $.DatasetID }}.{{ $t2 }}
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
	DatasetID           string
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
