package bigquery

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	dbOpt := dbutil.DBOpt{
		Backend:    types.BackendBigQuery,
		BigQueryDB: db.Client,
		DatasetID:  &db.datasetID,
	}
	doJoinOpt := sqlutil.DoJoinOpt{
		JoinOpt:             opt,
		QueryResults:        bigqueryQueryResults,
		ReadJoinResultQuery: READ_JOIN_RESULT_QUERY,
	}
	return sqlutil.DoJoin(ctx, dbOpt, doJoinOpt)
}

func bigqueryQueryResults(ctx context.Context, dbOpt dbutil.DBOpt, query string, header, tableNames []string) (*types.JoinResult, error) {
	rows, err := dbOpt.BigQueryDB.Query(query).Read(ctx)
	if err != nil {
		return nil, err
	}

	data := make(chan []interface{})
	var scanErr, dropErr error

	go func() {
		defer func() {
			if err = dropTemporaryTables(ctx, dbOpt.BigQueryDB, tableNames); err != nil {
				dropErr = err
			}
			defer close(data)
		}()
		for {
			recordMap := make(map[string]bigquery.Value)
			err = rows.Next(&recordMap)
			if err == iterator.Done {
				break
			}
			if err != nil {
				scanErr = err
			}
			record := make([]interface{}, 0, len(recordMap))
			for _, h := range header {
				record = append(record, recordMap[h])
			}
			data <- record
		}
	}()

	// TODO: return errors through channel
	if scanErr != nil {
		return nil, scanErr
	}

	return &types.JoinResult{
		Header: header,
		Data:   data,
	}, dropErr
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
