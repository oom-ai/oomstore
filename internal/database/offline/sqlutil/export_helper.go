package sqlutil

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/offline"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/pkg/errors"
)

const UNION_ENTITY_QUERY = `
INSERT INTO {{ qt .TableName }}
{{ snapshot .SnapshotTables }}
{{ cdc .CdcTables }}
`

const EXPORT_QUERY = `
SELECT
	e.{{ qt .EntityName }},
	{{ fieldJoin .Fields }}
FROM {{ qt .EntityTableName }} AS e
{{ range $table := .SnapshotTables }}
LEFT JOIN {{ qt $table }}
ON e.{{ qt $.EntityName }} = {{ qt $table }}.{{ qt $.EntityName }}
{{end}}
{{ range $table := .CdcTables }}
LEFT JOIN
(
	SELECT
		{{ $table }}_1.*
	FROM {{ qt $table }} AS {{ $table }}_1
	JOIN
	(SELECT
		{{ qt $.EntityName }},
		MAX({{ qt $.UnixMilli }}) AS {{ qt $.UnixMilli }}
	FROM {{ qt $table }}
	WHERE {{ qt $table }}.{{ qt $.UnixMilli }} <= ?
	GROUP BY {{ qt $.EntityName }}
	) AS {{ $table }}_2
	ON {{ $table }}_1.{{ qt $.EntityName }} = {{ $table }}_2.{{ qt $.EntityName }} AND {{ $table }}_1.{{ qt $.UnixMilli }} = {{ $table }}_2.{{ qt $.UnixMilli }}
	WHERE {{ $table }}_1.{{ qt $.UnixMilli }} <= ?
) AS {{ $table }}_0
ON e.{{ qt $.EntityName }} = {{ $table }}_0.{{ qt $.EntityName }}
{{end}}

`

type unionEntityQueryParams struct {
	TableName      string
	EntityName     string
	SnapshotTables []string
	CdcTables      []string
	UnixMilli      int64
	Backend        types.BackendType
}

func buildUnionEntityQuery(params unionEntityQueryParams) (string, []interface{}, error) {
	qt := dbutil.QuoteFn(params.Backend)
	var args []interface{}
	t := template.Must(template.New("union_entity").Funcs(template.FuncMap{
		"qt": qt,
		"snapshot": func(tables []string) string {
			query := make([]string, 0, len(tables))
			for _, t := range tables {
				query = append(query, fmt.Sprintf("SELECT %s FROM %s", qt(params.EntityName), qt(t)))
			}
			return strings.Join(query, "UNION \n\t")
		},
		"cdc": func(tables []string) string {
			if len(tables) == 0 {
				return ""
			}
			query := make([]string, 0, len(tables))
			for _, t := range tables {
				query = append(query, fmt.Sprintf("SELECT %s FROM %s WHERE %s <= ?", qt(params.EntityName), qt(t), qt("unix_milli")))
				args = append(args, params.UnixMilli)
			}
			return fmt.Sprintf("%s%s", "UNION \n\t", strings.Join(query, "UNION \n\t"))
		},
	}).Parse(UNION_ENTITY_QUERY))

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, params); err != nil {
		return "", nil, errors.WithStack(err)
	}
	return buf.String(), args, nil
}

type exportQueryParams struct {
	EntityTableName string
	EntityName      string
	UnixMilli       string
	SnapshotTables  []string
	CdcTables       []string
	Fields          []string
	Backend         types.BackendType
}

func buildExportQuery(params exportQueryParams) (string, error) {
	t := template.Must(template.New("export").Funcs(template.FuncMap{
		"qt": dbutil.QuoteFn(params.Backend),
		"fieldJoin": func(fields []string) string {
			return strings.Join(fields, ",\n\t")
		},
	}).Parse(EXPORT_QUERY))

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, params); err != nil {
		return "", errors.WithStack(err)
	}
	return buf.String(), nil
}

func prepareEntityTable(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.ExportOpt, snapshotTables, cdcTables []string) (string, error) {
	// Step 1: create table export_entity
	tableName := dbutil.TempTable("export_entity")
	qtTableName, columnDefs, err := prepareTableSchema(dbOpt, prepareTableSchemaParams{
		tableName:    tableName,
		entityName:   opt.EntityName,
		hasUnixMilli: false,
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

	// Step 2: aggregate all entity keys
	query, args, err := buildUnionEntityQuery(unionEntityQueryParams{
		TableName:      tableName,
		EntityName:     opt.EntityName,
		SnapshotTables: snapshotTables,
		CdcTables:      cdcTables,
		UnixMilli:      opt.UnixMilli,
		Backend:        dbOpt.Backend,
	})
	if err != nil {
		return "", errdefs.WithStack(err)
	}
	if err = dbOpt.ExecContext(ctx, query, args); err != nil {
		return "", errdefs.WithStack(err)
	}

	// Step 3: create index on table entity_rows
	if supportIndex(dbOpt.Backend) {
		qt := dbutil.QuoteFn(dbOpt.Backend)
		index := fmt.Sprintf(`CREATE UNIQUE INDEX idx_%s ON %s (%s)`, tableName, tableName, qt(opt.EntityName))
		if err = dbOpt.ExecContext(ctx, index, nil); err != nil {
			return "", errdefs.WithStack(err)
		}
	}
	return tableName, nil
}
