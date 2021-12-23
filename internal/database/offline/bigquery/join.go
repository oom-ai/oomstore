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
	// Step 1: prepare temporary table entity_rows
	features := types.FeatureList{}
	for _, featureList := range opt.FeatureMap {
		features = append(features, featureList...)
	}
	if len(features) == 0 {
		return nil, nil
	}
	dbOpt := dbutil.DBOpt{
		Backend:    types.BackendBigQuery,
		BigQueryDB: db.Client,
		DatasetID:  &db.datasetID,
	}
	entityRowsTableName, err := sqlutil.PrepareEntityRowsTable(ctx, dbOpt, opt.Entity, opt.EntityRows, opt.ValueNames)
	if err != nil {
		return nil, err
	}

	// Step 2: process features by group, insert result to table joined
	tableNames := make([]string, 0)
	tableToFeatureMap := make(map[string]types.FeatureList)
	for groupName, featureList := range opt.FeatureMap {
		revisionRanges, ok := opt.RevisionRangeMap[groupName]
		if !ok {
			continue
		}
		joinedTableName, err := sqlutil.JoinOneGroup(ctx, dbOpt, offline.JoinOneGroupOpt{
			GroupName:           groupName,
			Entity:              opt.Entity,
			Features:            featureList,
			RevisionRanges:      revisionRanges,
			EntityRowsTableName: entityRowsTableName,
		})
		if err != nil {
			return nil, err
		}
		if joinedTableName != "" {
			tableNames = append(tableNames, joinedTableName)
			tableToFeatureMap[joinedTableName] = featureList
		}
	}

	//// Step 3: read joined results
	return sqlutil.ReadJoinedTable(ctx, dbOpt, sqlutil.ReadJoinedTableOpt{
		EntityRowsTableName: entityRowsTableName,
		TableNames:          tableNames,
		FeatureMap:          tableToFeatureMap,
		ValueNames:          opt.ValueNames,
		QueryResults:        bigqueryQueryResults,
		ReadJoinResultQuery: READ_JOIN_RESULT_QUERY,
	})
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
