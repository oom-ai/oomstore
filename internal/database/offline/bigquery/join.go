package bigquery

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	dbOpt := dbutil.DBOpt{
		Backend:    Backend,
		BigQueryDB: db.Client,
		DatasetID:  &db.datasetID,
	}
	doJoinOpt := sqlutil.DoJoinOpt{
		JoinOpt:             opt,
		QueryResults:        bigqueryQueryResults,
		QueryTableTimeRange: bigqueryQueryTimeRange,
		ReadJoinResultQuery: READ_JOIN_RESULT_QUERY,
	}
	return sqlutil.DoJoin(ctx, dbOpt, doJoinOpt)
}

func bigqueryQueryTimeRange(ctx context.Context, dbOpt dbutil.DBOpt, tableName string) (*types.DataTableTimeRange, error) {
	return getTableTimeRange(ctx, dbOpt, tableName)
}

func bigqueryQueryResults(ctx context.Context, dbOpt dbutil.DBOpt, query string, header dbutil.ColumnList, dropTableNames []string, backendType types.BackendType) (*types.JoinResult, error) {
	rows, err := dbOpt.BigQueryDB.Query(query).Read(ctx)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}

	data := make(chan types.JoinRecord)
	go func() {
		defer func() {
			if err = dropTemporaryTables(ctx, dbOpt.BigQueryDB, dropTableNames); err != nil {
				select {
				case data <- types.JoinRecord{Error: err}:
					// nothing to do
				default:
				}
			}
			close(data)
		}()

		for {
			recordMap := make(map[string]bigquery.Value)
			err = rows.Next(&recordMap)
			if err == iterator.Done {
				return
			}
			if err != nil {
				select {
				case data <- types.JoinRecord{Error: errdefs.WithStack(err)}:
					// nothing to do
				case <-ctx.Done():
					return
				}
				continue
			}
			record := make([]interface{}, 0, len(recordMap))
			for i := range header {
				column := strings.Split(header[i].Name, ".")
				deserializedValue, err := dbutil.DeserializeByValueType(recordMap[column[len(column)-1]], header[i].ValueType, backendType)
				if err != nil {
					select {
					case data <- types.JoinRecord{Error: err}:
						// nothing to do
					case <-ctx.Done():
						return
					}
					continue
				}
				record = append(record, deserializedValue)
			}

			select {
			case data <- types.JoinRecord{Record: record, Error: nil}:
				// nothing to do
			case <-ctx.Done():
				return
			}
		}
	}()

	return &types.JoinResult{
		Header: header.Names(),
		Data:   data,
	}, nil
}

func dropTemporaryTables(ctx context.Context, db *bigquery.Client, tableNames []string) error {
	var err error
	for _, tableName := range tableNames {
		if tmpErr := dropTable(ctx, db, tableName); tmpErr != nil {
			err = tmpErr
		}
	}
	return err
}

func dropTable(ctx context.Context, db *bigquery.Client, tableName string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName)
	if _, err := db.Query(query).Read(ctx); err != nil {
		return errdefs.WithStack(err)
	}
	return nil
}

const READ_JOIN_RESULT_QUERY = `
SELECT
	{{ qt .EntityRowsTableName }}.{{ qt .EntityKey }},
	{{ qt .EntityRowsTableName }}.{{ qt .UnixMilli }},
	{{ fieldJoin .Fields }}
FROM {{ $.DatasetID }}.{{ qt .EntityRowsTableName }}
{{ range $pair := .JoinTables }}
	{{- $t1 := qt $pair.LeftTable -}}
	{{- $t2 := qt $pair.RightTable -}}
lEFT JOIN {{ $.DatasetID }}.{{ $t2 }}
ON {{ $t1 }}.{{ qt $.UnixMilli }} = {{ $t2 }}.{{ qt $.UnixMilli }} AND {{ $t1 }}.{{ qt $.EntityKey }} = {{ $t2 }}.{{ qt $.EntityKey }}
{{end}}`
