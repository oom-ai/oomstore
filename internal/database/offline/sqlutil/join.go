package sqlutil

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type QueryResults func(ctx context.Context, dbOpt dbutil.DBOpt, query string, header dbutil.ColumnList, dropTableNames []string, backendType types.BackendType) (*types.JoinResult, error)

type ReadJoinedTableOpt struct {
	EntityRowsTableName string
	TableNames          []string
	AllTableNames       []string
	FeatureMap          map[string]types.FeatureList
	ValueNames          []string
	QueryResults        QueryResults
	ReadJoinResultQuery string
	BackendType         types.BackendType
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
	data := make(chan []interface{})
	defer close(data)
	emptyResult := &types.JoinResult{
		Data: data,
	}
	// Step 1: prepare temporary table entity_rows
	features := types.FeatureList{}
	for _, featureList := range opt.FeatureMap {
		features = append(features, featureList...)
	}
	if len(features) == 0 {
		return emptyResult, nil
	}
	entityRowsTableName, err := PrepareEntityRowsTable(ctx, dbOpt, opt.Entity, opt.EntityRows, opt.ValueNames)
	if err != nil {
		return nil, err
	}

	// Step 2: process features by group, insert result to table joined
	tableNames := make([]string, 0)
	allTableNames := make([]string, 0)
	tableToFeatureMap := make(map[string]types.FeatureList)
	var category types.Category
	for groupName, featureList := range opt.FeatureMap {
		revisionRanges, ok := opt.RevisionRangeMap[groupName]
		if !ok {
			continue
		}
		if revisionRanges[0].CdcTable == "" {
			category = types.CategoryBatch
		} else {
			category = types.CategoryStream
		}
		joinedTables, err := JoinOneGroup(ctx, dbOpt, offline.JoinOneGroupOpt{
			GroupName:           groupName,
			Category:            category,
			Entity:              opt.Entity,
			Features:            featureList,
			RevisionRanges:      revisionRanges,
			EntityRowsTableName: entityRowsTableName,
		})
		if err != nil {
			return nil, err
		}
		allTableNames = append(allTableNames, joinedTables...)
		if len(joinedTables) != 0 {
			tableName := joinedTables[len(joinedTables)-1]
			tableNames = append(tableNames, tableName)
			tableToFeatureMap[tableName] = featureList
		}
	}

	// Step 3: read joined results
	return ReadJoinedTable(ctx, dbOpt, ReadJoinedTableOpt{
		EntityRowsTableName: entityRowsTableName,
		TableNames:          tableNames,
		AllTableNames:       allTableNames,
		FeatureMap:          tableToFeatureMap,
		ValueNames:          opt.ValueNames,
		QueryResults:        opt.QueryResults,
		ReadJoinResultQuery: opt.ReadJoinResultQuery,
		BackendType:         dbOpt.Backend,
	})
}

func JoinOneGroup(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.JoinOneGroupOpt) ([]string, error) {
	if len(opt.Features) == 0 {
		return nil, nil
	}

	// Step 1: create temporary joined table
	snapshotJoinedTableName, err := PrepareJoinedTable(ctx, dbOpt, opt.Features, opt.Entity, opt.GroupName, opt.ValueNames)
	if err != nil {
		return nil, err
	}

	// Step 2: iterate each table range, join entity_rows table and each data tables
	columns := append(opt.ValueNames, opt.Features.Names()...)
	for _, r := range opt.RevisionRanges {
		query, err := buildJoinQuery(joinQueryParams{
			TableName:           snapshotJoinedTableName,
			EntityName:          opt.Entity.Name,
			EntityKey:           "entity_key",
			UnixMilli:           "unix_milli",
			Columns:             columns,
			EntityRowsTableName: opt.EntityRowsTableName,
			SnapshotTable:       r.SnapshotTable,
			Backend:             dbOpt.Backend,
			DatasetID:           dbOpt.DatasetID,
		})
		if err != nil {
			return nil, err
		}
		if err = dbOpt.ExecContext(ctx, query, []interface{}{r.MinRevision, r.MaxRevision}); err != nil {
			return nil, err
		}
	}
	if opt.Category == types.CategoryBatch {
		return []string{snapshotJoinedTableName}, nil
	}

	// Step 3: for streaming features, keep joining with cdc_table
	cdcJoinedTableName, err := PrepareJoinedTable(ctx, dbOpt, opt.Features, opt.Entity, opt.GroupName, opt.ValueNames)
	if err != nil {
		return nil, err
	}
	for _, r := range opt.RevisionRanges {
		query, err := buildCdcJoinQuery(cdcJoinQueryParams{
			TableName:           cdcJoinedTableName,
			EntityKey:           "entity_key",
			EntityName:          opt.Entity.Name,
			UnixMilli:           "unix_milli",
			ValueNames:          opt.ValueNames,
			FeatureNames:        opt.Features.Names(),
			SnapshotJoinedTable: snapshotJoinedTableName,
			CdcTable:            r.CdcTable,
			Backend:             dbOpt.Backend,
			DatasetID:           dbOpt.DatasetID,
		})
		if err != nil {
			return nil, err
		}
		if err = dbOpt.ExecContext(ctx, query, []interface{}{r.MinRevision, r.MaxRevision}); err != nil {
			return nil, err
		}
	}
	return []string{snapshotJoinedTableName, cdcJoinedTableName}, nil
}

func ReadJoinedTable(ctx context.Context, dbOpt dbutil.DBOpt, opt ReadJoinedTableOpt) (*types.JoinResult, error) {
	data := make(chan []interface{})
	defer close(data)
	emptyResult := &types.JoinResult{
		Data: data,
	}
	tableNames := opt.TableNames
	dropTableNames := append([]string{opt.EntityRowsTableName}, opt.AllTableNames...)
	if len(tableNames) == 0 {
		return emptyResult, nil
	}
	qt := dbutil.QuoteFn(dbOpt.Backend)

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
		fields         []string
		header         []dbutil.Column
		joinTablePairs []joinTablePair
	)
	header = append(header, dbutil.Column{
		Name:      "entity_key",
		ValueType: types.String,
	}, dbutil.Column{
		Name:      "unix_milli",
		ValueType: types.Int64,
	})
	for _, name := range opt.ValueNames {
		fields = append(fields, qt(name))
		header = append(header, dbutil.Column{
			Name:      name,
			ValueType: types.String,
		})
	}
	for _, name := range tableNames {
		for _, f := range opt.FeatureMap[name] {
			fields = append(fields, fmt.Sprintf("%s.%s", qt(name), qt(f.Name)))
			header = append(header, dbutil.Column{
				Name:      f.FullName,
				ValueType: f.ValueType,
			})
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
		EntityKey:           "entity_key",
		UnixMilli:           "unix_milli",
		Fields:              fields,
		JoinTables:          joinTablePairs,
		Backend:             dbOpt.Backend,
		DatasetID:           datasetID,
	})
	if err != nil {
		return nil, err
	}

	// Step 2: read joined results
	return opt.QueryResults(ctx, dbOpt, query, header, dropTableNames, dbOpt.Backend)
}

func sqlxQueryResults(ctx context.Context, dbOpt dbutil.DBOpt, query string, header dbutil.ColumnList, dropTableNames []string, backendType types.BackendType) (*types.JoinResult, error) {
	stmt, err := dbOpt.SqlxDB.Preparex(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Queryx()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	data := make(chan []interface{})
	var scanErr, dropErr error
	go func() {
		defer func() {
			if err := dropTemporaryTables(ctx, dbOpt.SqlxDB, dropTableNames); err != nil {
				dropErr = err
			}
			defer rows.Close()
			defer close(data)
		}()
		for rows.Next() {
			record, err := rows.SliceScan()
			if err != nil {
				scanErr = errors.WithStack(err)
				continue
			}
			serializedRecord := make([]interface{}, 0, len(record))
			for i, r := range record {
				deserializedValue, err := dbutil.DeserializeByValueType(r, header[i].ValueType, backendType)
				if err != nil {
					scanErr = err
					continue
				}
				serializedRecord = append(serializedRecord, deserializedValue)
			}
			data <- serializedRecord
		}
	}()

	// TODO: return errors through channel
	if scanErr != nil {
		return nil, scanErr
	}
	return &types.JoinResult{
		Header: header.Names(),
		Data:   data,
	}, dropErr
}
