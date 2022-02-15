package sqlutil

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type QueryResults func(ctx context.Context, dbOpt dbutil.DBOpt, query string, header dbutil.ColumnList, dropTableNames []string, backendType types.BackendType) (*types.JoinResult, error)
type QueryTableTimeRange func(ctx context.Context, dbOpt dbutil.DBOpt, tableName string) (*types.DataTableTimeRange, error)

func Join(ctx context.Context, db *sqlx.DB, opt offline.JoinOpt, backendType types.BackendType) (*types.JoinResult, error) {
	dbOpt := dbutil.DBOpt{
		Backend: backendType,
		SqlxDB:  db,
	}
	doJoinOpt := DoJoinOpt{
		JoinOpt:             opt,
		QueryResults:        sqlxQueryResults,
		QueryTableTimeRange: sqlxQueryTableTimeRange,
		ReadJoinResultQuery: READ_JOIN_RESULT_QUERY,
	}
	return DoJoin(ctx, dbOpt, doJoinOpt)
}

type DoJoinOpt struct {
	offline.JoinOpt
	QueryResults        QueryResults
	QueryTableTimeRange QueryTableTimeRange
	ReadJoinResultQuery string
}

func DoJoin(ctx context.Context, dbOpt dbutil.DBOpt, opt DoJoinOpt) (*types.JoinResult, error) {
	if err := validateJoinOpt(opt); err != nil {
		data := make(chan types.JoinRecord)
		close(data)
		return &types.JoinResult{Data: data}, nil
	}

	// Step 1: prepare temporary table entity_rows
	entityRowsTableName, err := prepareEntityRowsTable(ctx, dbOpt, opt.EntityRows, opt.ValueNames)
	if err != nil {
		return nil, err
	}
	timeRange, err := opt.QueryTableTimeRange(ctx, dbOpt, entityRowsTableName)
	if err != nil {
		return nil, err
	}
	if timeRange.MinUnixMilli == nil || timeRange.MaxUnixMilli == nil {
		data := make(chan types.JoinRecord)
		close(data)
		return &types.JoinResult{Data: data}, nil
	}

	// Step 2: process features by group, insert result to table joined
	tableNames := make([]string, 0)
	allTableNames := make([]string, 0)
	tableToFeatureMap := make(map[string]types.FeatureList)
	var category types.Category
	for _, groupName := range opt.GroupNames {
		featureList := opt.FeatureMap[groupName]
		revisionRanges, ok := opt.RevisionRangeMap[groupName]
		if !ok {
			continue
		}
		if revisionRanges[0].CdcTable == "" {
			category = types.CategoryBatch
		} else {
			category = types.CategoryStream
		}
		joinedTables, err := joinOneGroup(ctx, dbOpt, joinOneGroupOpt{
			GroupName:           groupName,
			Category:            category,
			EntityName:          opt.EntityName,
			Features:            featureList,
			RevisionRanges:      revisionRanges,
			EntityRowsTableName: entityRowsTableName,
			TimeRange:           *timeRange,
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
	return readJoinedTable(ctx, dbOpt, readJoinedTableOpt{
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

func validateJoinOpt(opt DoJoinOpt) error {
	for groupName, features := range opt.FeatureMap {
		if len(features) == 0 {
			delete(opt.FeatureMap, groupName)
			delete(opt.RevisionRangeMap, groupName)
		}
	}
	if len(opt.FeatureMap) == 0 || len(opt.RevisionRangeMap) == 0 {
		return errdefs.Errorf("empty feature map")
	}
	return nil
}

type joinOneGroupOpt struct {
	GroupName           string
	Category            types.Category
	Features            types.FeatureList
	RevisionRanges      []*offline.RevisionRange
	EntityName          string
	EntityRowsTableName string
	ValueNames          []string
	TimeRange           types.DataTableTimeRange
}

func joinOneGroup(ctx context.Context, dbOpt dbutil.DBOpt, opt joinOneGroupOpt) ([]string, error) {
	if len(opt.Features) == 0 {
		return nil, nil
	}

	// Step 1: create temporary joined table
	snapshotJoinedTableName, err := prepareJoinedTable(ctx, dbOpt, opt.Features, opt.GroupName, opt.ValueNames)
	if err != nil {
		return nil, err
	}

	// Step 2: iterate each table range, join entity_rows table and each data tables
	columns := append(opt.ValueNames, opt.Features.Names()...)
	for _, r := range opt.RevisionRanges {
		if *opt.TimeRange.MaxUnixMilli < r.MinRevision || *opt.TimeRange.MinUnixMilli > r.MaxRevision {
			continue
		}
		query, err := buildJoinQuery(joinQueryParams{
			TableName:           snapshotJoinedTableName,
			EntityName:          opt.EntityName,
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
	cdcJoinedTableName, err := prepareJoinedTable(ctx, dbOpt, opt.Features, opt.GroupName, opt.ValueNames)
	if err != nil {
		return nil, err
	}
	for _, r := range opt.RevisionRanges {
		query, err := buildCdcJoinQuery(cdcJoinQueryParams{
			TableName:           cdcJoinedTableName,
			EntityKey:           "entity_key",
			EntityName:          opt.EntityName,
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

type readJoinedTableOpt struct {
	EntityRowsTableName string
	TableNames          []string
	AllTableNames       []string
	FeatureMap          map[string]types.FeatureList
	ValueNames          []string
	QueryResults        QueryResults
	ReadJoinResultQuery string
	BackendType         types.BackendType
}

func readJoinedTable(ctx context.Context, dbOpt dbutil.DBOpt, opt readJoinedTableOpt) (*types.JoinResult, error) {
	if len(opt.TableNames) == 0 {
		data := make(chan types.JoinRecord)
		close(data)
		return &types.JoinResult{Data: data}, nil
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
	for _, name := range opt.TableNames {
		for _, f := range opt.FeatureMap[name] {
			fields = append(fields, fmt.Sprintf("%s.%s", qt(name), qt(f.Name)))
			header = append(header, dbutil.Column{
				Name:      f.FullName(),
				ValueType: f.ValueType,
			})
		}
	}

	tableNames := append([]string{opt.EntityRowsTableName}, opt.TableNames...)
	for i := 0; i < len(tableNames)-1; i++ {
		joinTablePairs = append(joinTablePairs, joinTablePair{
			LeftTable:  tableNames[0],
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
	dropTableNames := []string{buildTableName(dbOpt, opt.EntityRowsTableName)}
	for _, tableName := range opt.AllTableNames {
		dropTableName := buildTableName(dbOpt, tableName)
		dropTableNames = append(dropTableNames, dropTableName)
		if err := AddTemporaryTableRecord(ctx, dbOpt, dropTableName); err != nil {
			return nil, err
		}
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
		return nil, errdefs.WithStack(err)
	}

	data := make(chan types.JoinRecord)
	go func() {
		defer func() {
			if err := DropTemporaryTables(ctx, dbOpt, dropTableNames); err != nil {
				select {
				case data <- types.JoinRecord{Error: err}:
					// nothing to do
				default:
				}
			}

			rows.Close()
			close(data)
		}()

		for rows.Next() {
			record, err := rows.SliceScan()
			if err != nil {
				data <- types.JoinRecord{Error: errdefs.WithStack(err)}
				continue
			}

			deserializedRecord := make([]interface{}, 0, len(record))
			for i, r := range record {
				deserializedValue, err := dbutil.DeserializeByValueType(r, header[i].ValueType, backendType)
				if err != nil {
					data <- types.JoinRecord{Error: err}
				}
				deserializedRecord = append(deserializedRecord, deserializedValue)
			}

			select {
			case data <- types.JoinRecord{Record: deserializedRecord, Error: nil}:
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

func sqlxQueryTableTimeRange(ctx context.Context, dbOpt dbutil.DBOpt, tableName string) (*types.DataTableTimeRange, error) {
	return getCdcTimeRange(ctx, dbOpt.SqlxDB, tableName, dbOpt.Backend)
}
