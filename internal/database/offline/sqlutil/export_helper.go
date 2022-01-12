package sqlutil

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/pkg/errors"
)

const AGGREGATE_QUERY = `
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
		WHERE {{ qt .UnixMilli }} <= ?
		GROUP BY {{ qt .EntityName }}) AS t2
	ON t1.{{ qt .EntityName }} = t2.{{ qt .EntityName }} AND t1.{{ qt .UnixMilli }} = t2.{{ qt .UnixMilli }}
	WHERE t1.{{ qt .UnixMilli }} <= ?
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
		WHERE {{ qt .UnixMilli }} <= ?
		GROUP BY {{ qt .EntityName }}) AS t2
	ON t1.{{ qt .EntityName }} = t2.{{ qt .EntityName }} AND t1.{{ qt .UnixMilli }} = t2.{{ qt .UnixMilli }}
	WHERE t1.{{ qt .UnixMilli }} <= ?
) AS r
LEFT JOIN
{{ qt .PrevSnapshotTableName }} AS l
ON l.{{ qt .EntityName }} = r.{{ qt .EntityName }}
`

type aggregateQueryParams struct {
	EntityName            string
	UnixMilli             string
	FeatureNames          []string
	PrevSnapshotTableName string
	CurrCdcTableName      string
	Union                 string
	Backend               types.BackendType
	DatasetID             *string
}

func buildAggregateQuery(params aggregateQueryParams) (string, error) {
	if params.Backend == types.BackendBigQuery {
		params.Union = "UNION DISTINCT"
		params.PrevSnapshotTableName = fmt.Sprintf("%s.%s", *params.DatasetID, params.PrevSnapshotTableName)
		params.CurrCdcTableName = fmt.Sprintf("%s.%s", *params.DatasetID, params.CurrCdcTableName)
	} else {
		params.Union = "UNION"
	}
	qt := dbutil.QuoteFn(params.Backend)
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
	}).Parse(AGGREGATE_QUERY))

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, params); err != nil {
		return "", errors.WithStack(err)
	}
	return buf.String(), nil
}
