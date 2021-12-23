package bigquery

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

func dropTable(ctx context.Context, db *bigquery.Client, tableName string) error {
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
