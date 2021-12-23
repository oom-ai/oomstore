package bigquery

import (
	"context"
	"fmt"

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
	return readJoinedTable(ctx, db, entityRowsTableName, tableNames, tableToFeatureMap, opt.ValueNames, types.BackendMySQL)
}

func readJoinedTable(
	ctx context.Context,
	db *DB,
	entityRowsTableName string,
	tableNames []string,
	featureMap map[string]types.FeatureList,
	valueNames []string,
	backendType types.BackendType,
) (*types.JoinResult, error) {
	if len(tableNames) == 0 {
		return nil, nil
	}
	qt, err := dbutil.QuoteFn(backendType)
	if err != nil {
		return nil, err
	}
	entityKeyStr := qt("entity_key")
	unixMilliStr := qt("unix_milli")

	// Step 1: join temporary tables
	/*
		SELECT
		entity_rows_table.entity_key,
			entity_rows_table.unix_milli,
			joined_table_1.feature_1,
			joined_table_1.feature_2,
			joined_table_2.feature_3
		FROM entity_rows_table
		LEFT JOIN joined_table_1
		ON entity_rows_table.entity_key = joined_table_1.entity_key AND entity_rows_table.unix_milli = joined_table_1.unix_milli
		LEFT JOIN joined_table_2
		ON entity_rows_table.entity_key = joined_table_2.entity_key AND entity_rows_table.unix_milli = joined_table_2.unix_milli;
	*/
	var (
		header         []string
		fields         []string
		joinTablePairs []joinTablePair
	)

	header = append(header, "entity_key", "unix_milli")
	for _, name := range valueNames {
		fields = append(fields, qt(name))
		header = append(header, name)
	}
	for _, name := range tableNames {
		for _, f := range featureMap[name] {
			fields = append(fields, fmt.Sprintf("%s.%s", qt(name), qt(f.Name)))
			header = append(header, f.Name)
		}
	}

	tableNames = append([]string{entityRowsTableName}, tableNames...)
	for i := 0; i < len(tableNames)-1; i++ {
		joinTablePairs = append(joinTablePairs, joinTablePair{
			LeftTable:  tableNames[i],
			RightTable: tableNames[i+1],
		})
	}
	query, err := buildReadJoinResultQuery(readJoinResultQuery{
		EntityRowsTableName: entityRowsTableName,
		EntityKeyStr:        entityKeyStr,
		UnixMilliStr:        unixMilliStr,
		Fields:              fields,
		JoinTables:          joinTablePairs,
		Backend:             backendType,
		DatasetID:           db.datasetID,
	})
	if err != nil {
		return nil, err
	}
	// Step 2: read joined results
	rows, err := db.Query(query).Read(ctx)
	if err != nil {
		return nil, err
	}

	data := make(chan []interface{})
	var scanErr, dropErr error

	go func() {
		defer func() {
			if err := dropTemporaryTables(ctx, db, tableNames); err != nil {
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

func dropTemporaryTables(ctx context.Context, db *DB, tableNames []string) error {
	var err error
	for _, tableName := range tableNames {
		if tmpErr := dropTable(ctx, db, tableName); tmpErr != nil {
			err = tmpErr
		}
	}
	return err
}
