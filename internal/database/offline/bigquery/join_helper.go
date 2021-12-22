package bigquery

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
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

func dropTable(ctx context.Context, db *DB, tableName string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName)
	_, err := db.Query(query).Read(ctx)
	return err
}

func convertValueTypeToBigQuerySQLType(t types.ValueType) (string, error) {
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
