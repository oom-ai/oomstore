package sqlutil

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"

	"github.com/oom-ai/oomstore/internal/database/metadata/sqlutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
)

const SNAPSHOT_QUERY = `
INSERT INTO {{ qt .CurrSnapshotTableName }} ({{ qt .EntityName }}, {{ columnJoin .FeatureNames }})
SELECT
	l.{{ qt .EntityName }},
	{{ featureValue .FeatureNames }}
FROM {{ qt .PrevSnapshotTableName }} AS l
LEFT JOIN
(
	SELECT
		t1.*
	FROM
		{{ qt .CurrCdcTableName }} AS t1
	JOIN
		(SELECT
			{{ qt .EntityName }},
			MAX({{ qt .UnixMilli }}) AS {{ qt .UnixMilli }}
		FROM {{ qt .CurrCdcTableName }}
		GROUP BY {{ qt .EntityName }}) AS t2
	ON t1.{{ qt .EntityName }} = t2.{{ qt .EntityName }} AND t1.{{ qt .UnixMilli }} = t2.{{ qt .UnixMilli }}
) AS r
ON l.{{ qt .EntityName }} = r.{{ qt .EntityName }}
{{ .Union }}
SELECT
	r.{{ qt .EntityName }},
	{{ featureValue .FeatureNames }}
FROM
(
	SELECT
		t1.*
	FROM
		{{ qt .CurrCdcTableName }} AS t1
	JOIN
		(SELECT
			{{ qt .EntityName }},
			MAX({{ qt .UnixMilli }}) AS {{ qt .UnixMilli }}
		FROM {{ qt .CurrCdcTableName }}
		GROUP BY {{ qt .EntityName }}) AS t2
	ON t1.{{ qt .EntityName }} = t2.{{ qt .EntityName }} AND t1.{{ qt .UnixMilli }} = t2.{{ qt .UnixMilli }}
) AS r
LEFT JOIN
{{ qt .PrevSnapshotTableName }} AS l
ON l.{{ qt .EntityName }} = r.{{ qt .EntityName }}
`

func Snapshot(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.SnapshotOpt) error {
	prevSnapshotTableName := sqlutil.OfflineStreamSnapshotTableName(opt.Group.ID, opt.PrevRevisionID)
	currSnapshotTableName := sqlutil.OfflineStreamSnapshotTableName(opt.Group.ID, opt.RevisionID)
	currCdcTableName := sqlutil.OfflineStreamCdcTableName(opt.Group.ID, opt.RevisionID)

	if dbOpt.Backend == types.BackendBigQuery {
		currSnapshotTableName = fmt.Sprintf("%s.%s", *dbOpt.DatasetID, currSnapshotTableName)
		prevSnapshotTableName = fmt.Sprintf("%s.%s", *dbOpt.DatasetID, prevSnapshotTableName)
		currCdcTableName = fmt.Sprintf("%s.%s", *dbOpt.DatasetID, currCdcTableName)
	}

	schema, err := dbutil.BuildCreateSchema(currSnapshotTableName, opt.Group.Entity, opt.Features, dbOpt.Backend)
	if err != nil {
		return err
	}

	if err = dbOpt.ExecContext(ctx, schema, nil); err != nil {
		return err
	}
	query, err := buildSnapshotQuery(snapshotQueryParams{
		EntityName:            opt.Group.Entity.Name,
		UnixMilli:             "unix_milli",
		FeatureNames:          opt.Features.Names(),
		PrevSnapshotTableName: prevSnapshotTableName,
		CurrSnapshotTableName: currSnapshotTableName,
		CurrCdcTableName:      currCdcTableName,
		Backend:               dbOpt.Backend,
	})
	if err != nil {
		return err
	}
	if err = dbOpt.ExecContext(ctx, query, nil); err != nil {
		return err
	}

	return nil
}

type snapshotQueryParams struct {
	EntityName            string
	UnixMilli             string
	FeatureNames          []string
	PrevSnapshotTableName string
	CurrSnapshotTableName string
	CurrCdcTableName      string
	Union                 string
	Backend               types.BackendType
	DatasetID             *string
}

func buildSnapshotQuery(params snapshotQueryParams) (string, error) {
	if params.Backend == types.BackendBigQuery {
		params.Union = "UNION DISTINCT"
	} else {
		params.Union = "UNION"
	}
	qt, err := dbutil.QuoteFn(params.Backend)
	if err != nil {
		return "", err
	}

	t := template.Must(template.New("snapshot").Funcs(template.FuncMap{
		"qt": qt,
		"columnJoin": func(columns []string) string {
			return qt(columns...)
		},
		"featureValue": func(features []string) string {
			values := make([]string, 0, len(features))
			for _, f := range features {
				values = append(values, fmt.Sprintf("(CASE WHEN r.%s IS NULL THEN l.%s ELSE r.%s END) AS %s", qt(f), qt(f), qt(f), f))
			}
			return strings.Join(values, ",")
		},
	}).Parse(SNAPSHOT_QUERY))

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}