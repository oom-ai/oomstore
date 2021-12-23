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

func dropTable(ctx context.Context, db *DB, tableName string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName)
	_, err := db.Query(query).Read(ctx)
	return err
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
