package postgres

import (
	"context"
	"fmt"
	"strings"

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
	entityRowsTableName, err := db.createAndImportTableEntityRows(ctx, opt.Entity, opt.EntityRows)
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
		joinedTableName, err := db.joinOneGroup(ctx, offline.JoinOneGroupOpt{
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
	return db.readJoinedTable(ctx, entityRowsTableName, tableNames, tableToFeatureMap)
}

func (db *DB) joinOneGroup(ctx context.Context, opt offline.JoinOneGroupOpt) (string, error) {
	if len(opt.Features) == 0 {
		return "", nil
	}
	// Step 1: create temporary joined table
	joinedTableName, err := db.createTableJoined(ctx, opt.Features, opt.Entity, opt.GroupName)
	if err != nil {
		return "", err
	}

	// Step 2: iterate each table range, join entity_rows table and each data tables
	joinQuery := `
		INSERT INTO "%s"(entity_key, unix_time, %s)
		SELECT
			l.entity_key AS entity_key,
			l.unix_time AS unix_time,
			%s
		FROM "%s" AS l
		LEFT JOIN "%s" AS r
		ON l.entity_key = r.%s
		WHERE l.unix_time >= $1 AND l.unix_time < $2;
	`
	featureNamesStr := dbutil.Quote(`"`, opt.Features.Names()...)
	for _, r := range opt.RevisionRanges {
		query := fmt.Sprintf(joinQuery, joinedTableName, featureNamesStr, featureNamesStr, opt.EntityRowsTableName, r.DataTable, opt.Entity.Name)
		if _, tmpErr := db.ExecContext(ctx, query, r.MinRevision, r.MaxRevision); tmpErr != nil {
			return "", tmpErr
		}
	}

	return joinedTableName, nil
}

func (db *DB) readJoinedTable(ctx context.Context, entityRowsTableName string, tableNames []string, featureMap map[string]types.FeatureList) (*types.JoinResult, error) {
	if len(tableNames) == 0 {
		return nil, nil
	}

	// Step 1: join temporary tables
	/*
		SELECT
		entity_rows_table.entity_key,
			entity_rows_table.unix_time,
			joined_table_1.feature_1,
			joined_table_1.feature_2,
			joined_table_2.feature_3
		FROM entity_rows_table
		LEFT JOIN joined_table_1
		ON entity_rows_table.entity_key = joined_table_1.entity_key AND entity_rows_table.unix_time = joined_table_1.unix_time
		LEFT JOIN joined_table_2
		ON entity_rows_table.entity_key = joined_table_2.entity_key AND entity_rows_table.unix_time = joined_table_2.unix_time;
	*/
	fields := []string{fmt.Sprintf("%s.entity_key, %s.unix_time", entityRowsTableName, entityRowsTableName)}
	for _, tableName := range tableNames {
		for _, f := range featureMap[tableName] {
			fields = append(fields, fmt.Sprintf("%s.%s", tableName, f.Name))
		}
	}
	query := fmt.Sprintf(`SELECT %s FROM %s`, strings.Join(fields, ","), dbutil.Quote(`"`, entityRowsTableName))
	tableNames = append([]string{entityRowsTableName}, tableNames...)
	for i := range tableNames {
		if i == 0 {
			continue
		}
		query = fmt.Sprintf("%s LEFT JOIN %s ON %s.unix_time = %s.unix_time AND %s.entity_key = %s.entity_key",
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
			if err := db.dropTemporaryTables(ctx, tableNames); err != nil {
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

func (db *DB) dropTemporaryTables(ctx context.Context, tableNames []string) error {
	var err error
	for _, tableName := range tableNames {
		if tmpErr := db.dropTable(ctx, tableName); tmpErr != nil {
			err = tmpErr
		}
	}
	return err
}
