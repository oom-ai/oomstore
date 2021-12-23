package sqlutil

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type QueryResults func(ctx context.Context, dbOpt dbutil.DBOpt, query string, header, tableNames []string) (*types.JoinResult, error)

type ReadJoinedTableOpt struct {
	EntityRowsTableName string
	TableNames          []string
	FeatureMap          map[string]types.FeatureList
	ValueNames          []string
	QueryResults        QueryResults
	ReadJoinResultQuery string
}

type DoJoinOpt struct {
	offline.JoinOpt
	QueryResults        QueryResults
	ReadJoinResultQuery string
}

func Join(ctx context.Context, db *sqlx.DB, opt offline.JoinOpt, backendType types.BackendType) (*types.JoinResult, error) {
	dbOpt := dbutil.DBOpt{
		Backend: backendType,
		SqlxDB:  db,
	}
	doJoinOpt := DoJoinOpt{
		JoinOpt:             opt,
		QueryResults:        sqlxQueryResults,
		ReadJoinResultQuery: READ_JOIN_RESULT_QUERY,
	}
	return DoJoin(ctx, dbOpt, doJoinOpt)
}

func DoJoin(ctx context.Context, dbOpt dbutil.DBOpt, opt DoJoinOpt) (*types.JoinResult, error) {
	// Step 1: prepare temporary table entity_rows
	features := types.FeatureList{}
	for _, featureList := range opt.FeatureMap {
		features = append(features, featureList...)
	}
	if len(features) == 0 {
		return nil, nil
	}
	entityRowsTableName, err := PrepareEntityRowsTable(ctx, dbOpt, opt.Entity, opt.EntityRows, opt.ValueNames)
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
		joinedTableName, err := JoinOneGroup(ctx, dbOpt, offline.JoinOneGroupOpt{
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

	// Step 3: read joined results
	return ReadJoinedTable(ctx, dbOpt, ReadJoinedTableOpt{
		EntityRowsTableName: entityRowsTableName,
		TableNames:          tableNames,
		FeatureMap:          tableToFeatureMap,
		ValueNames:          opt.ValueNames,
		QueryResults:        opt.QueryResults,
		ReadJoinResultQuery: opt.ReadJoinResultQuery,
	})
}

func JoinOneGroup(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.JoinOneGroupOpt) (string, error) {
	if len(opt.Features) == 0 {
		return "", nil
	}
	qt, err := dbutil.QuoteFn(dbOpt.Backend)
	if err != nil {
		return "", err
	}
	entityKeyStr := qt("entity_key")
	unixMilliStr := qt("unix_milli")

	// Step 1: create temporary joined table
	joinedTableName, err := PrepareJoinedTable(ctx, dbOpt, opt.Features, opt.Entity, opt.GroupName, opt.ValueNames)
	if err != nil {
		return "", err
	}

	// Step 2: iterate each table range, join entity_rows table and each data tables
	columns := append(opt.ValueNames, opt.Features.Names()...)
	for _, r := range opt.RevisionRanges {
		query, err := buildJoinQuery(joinQueryParams{
			TableName:           joinedTableName,
			EntityKeyStr:        entityKeyStr,
			EntityName:          opt.Entity.Name,
			UnixMilliStr:        unixMilliStr,
			Columns:             columns,
			EntityRowsTableName: opt.EntityRowsTableName,
			DataTable:           r.DataTable,
			Backend:             dbOpt.Backend,
			DatasetID:           dbOpt.DatasetID,
		})
		if err != nil {
			return "", err
		}
		if err = dbOpt.ExecContext(ctx, query, []interface{}{r.MinRevision, r.MaxRevision}); err != nil {
			return "", err
		}
	}

	return joinedTableName, nil
}

func ReadJoinedTable(ctx context.Context, dbOpt dbutil.DBOpt, opt ReadJoinedTableOpt) (*types.JoinResult, error) {
	tableNames := opt.TableNames
	if len(tableNames) == 0 {
		return nil, nil
	}
	qt, err := dbutil.QuoteFn(dbOpt.Backend)
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
		fields, header []string
		joinTablePairs []joinTablePair
	)
	header = append(header, "entity_key", "unix_milli")
	for _, name := range opt.ValueNames {
		fields = append(fields, qt(name))
		header = append(header, name)
	}
	for _, name := range tableNames {
		for _, f := range opt.FeatureMap[name] {
			fields = append(fields, fmt.Sprintf("%s.%s", qt(name), qt(f.Name)))
			header = append(header, f.Name)
		}
	}

	tableNames = append([]string{opt.EntityRowsTableName}, tableNames...)
	for i := 0; i < len(tableNames)-1; i++ {
		joinTablePairs = append(joinTablePairs, joinTablePair{
			LeftTable:  tableNames[i],
			RightTable: tableNames[i+1],
		})
	}
	datasetID := ""
	if dbOpt.Backend == types.BackendBigQuery {
		datasetID = *dbOpt.DatasetID
	}
	query, err := buildReadJoinResultQuery(opt.ReadJoinResultQuery, readJoinResultQueryParams{
		EntityRowsTableName: opt.EntityRowsTableName,
		EntityKeyStr:        entityKeyStr,
		UnixMilliStr:        unixMilliStr,
		Fields:              fields,
		JoinTables:          joinTablePairs,
		Backend:             dbOpt.Backend,
		DatasetID:           datasetID,
	})
	if err != nil {
		return nil, err
	}

	// Step 2: read joined results
	return opt.QueryResults(ctx, dbOpt, query, header, tableNames)
}

func sqlxQueryResults(ctx context.Context, dbOpt dbutil.DBOpt, query string, header, tableNames []string) (*types.JoinResult, error) {
	rows, err := dbOpt.SqlxDB.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	data := make(chan []interface{})
	var scanErr, dropErr error
	go func() {
		defer func() {
			if err := dropTemporaryTables(ctx, dbOpt.SqlxDB, tableNames); err != nil {
				dropErr = err
			}
			defer rows.Close()
			defer close(data)
		}()
		for rows.Next() {
			record, err := rows.SliceScan()
			if err != nil {
				scanErr = err
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
