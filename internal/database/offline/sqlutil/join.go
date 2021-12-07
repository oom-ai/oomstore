package sqlutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
)

func Join(ctx context.Context, db *sqlx.DB, opt offline.JoinOpt, backendType types.BackendType) (*types.JoinResult, error) {
	// Step 1: prepare temporary table entity_rows
	features := types.FeatureList{}
	for _, featureList := range opt.FeatureMap {
		features = append(features, featureList...)
	}
	if len(features) == 0 {
		return nil, nil
	}
	entityRowsTableName, err := createAndImportTableEntityRows(ctx, db, opt.Entity, opt.EntityRows, opt.ValueNames, backendType)
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
		joinedTableName, err := joinOneGroup(ctx, db, offline.JoinOneGroupOpt{
			GroupName:           groupName,
			Entity:              opt.Entity,
			Features:            featureList,
			RevisionRanges:      revisionRanges,
			EntityRowsTableName: entityRowsTableName,
		}, backendType)
		if err != nil {
			return nil, err
		}
		if joinedTableName != "" {
			tableNames = append(tableNames, joinedTableName)
			tableToFeatureMap[joinedTableName] = featureList
		}
	}

	// Step 3: read joined results
	return readJoinedTable(ctx, db, entityRowsTableName, tableNames, tableToFeatureMap, opt.ValueNames, backendType)
}

func joinOneGroup(ctx context.Context, db *sqlx.DB, opt offline.JoinOneGroupOpt, backendType types.BackendType) (string, error) {
	if len(opt.Features) == 0 {
		return "", nil
	}
	qt, err := dbutil.QuoteFn(backendType)
	if err != nil {
		return "", err
	}

	// Step 1: create temporary joined table
	joinedTableName, err := createTableJoined(ctx, db, opt.Features, opt.Entity, opt.GroupName, opt.ValueNames, backendType)
	if err != nil {
		return "", err
	}

	// Step 2: iterate each table range, join entity_rows table and each data tables
	joinQuery := `
		INSERT INTO %s(entity_key, unix_milli, %s)
		SELECT
			l.entity_key AS entity_key,
			l.unix_milli AS unix_milli,
			%s
		FROM %s AS l
		LEFT JOIN %s AS r
		ON l.entity_key = r.%s
		WHERE l.unix_milli >= ? AND l.unix_milli < ?;
	`

	columns := append(opt.ValueNames, opt.Features.Names()...)
	columnsStr := qt(columns...)
	for _, r := range opt.RevisionRanges {
		query := fmt.Sprintf(joinQuery,
			qt(joinedTableName),
			columnsStr,
			columnsStr,
			qt(opt.EntityRowsTableName),
			qt(r.DataTable),
			opt.Entity.Name,
		)
		if _, tmpErr := db.ExecContext(ctx, db.Rebind(query), r.MinRevision, r.MaxRevision); tmpErr != nil {
			return "", tmpErr
		}
	}

	return joinedTableName, nil
}

func readJoinedTable(
	ctx context.Context,
	db *sqlx.DB,
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
	fields := []string{fmt.Sprintf("%s.entity_key, %s.unix_milli", entityRowsTableName, entityRowsTableName)}
	for _, valueName := range valueNames {
		fields = append(fields, fmt.Sprintf("%s.%s", entityRowsTableName, valueName))
	}
	for _, tableName := range tableNames {
		for _, f := range featureMap[tableName] {
			fields = append(fields, fmt.Sprintf("%s.%s", tableName, f.Name))
		}
	}
	query := fmt.Sprintf(`SELECT %s FROM %s`, strings.Join(fields, ","), qt(entityRowsTableName))
	tableNames = append([]string{entityRowsTableName}, tableNames...)
	for i := range tableNames {
		if i == 0 {
			continue
		}
		query = fmt.Sprintf("%s LEFT JOIN %s ON %s.unix_milli = %s.unix_milli AND %s.entity_key = %s.entity_key",
			query, tableNames[i], tableNames[i-1], tableNames[i], tableNames[i-1], tableNames[i])
	}

	// Step 2: read joined results
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	header, err := rows.Columns()
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

func dropTemporaryTables(ctx context.Context, db *sqlx.DB, tableNames []string) error {
	var err error
	for _, tableName := range tableNames {
		if tmpErr := dropTable(ctx, db, tableName); tmpErr != nil {
			err = tmpErr
		}
	}
	return err
}
