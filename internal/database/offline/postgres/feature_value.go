package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (dataMap map[string]dbutil.RowMap, err error) {
	if len(opt.Features) == 0 {
		return make(map[string]dbutil.RowMap), nil
	}

	// Step 0: prepare temporary tables
	entityDfWithFeatureName, tmpErr := db.createTableEntityDfWithFeatures(ctx, opt.Features, opt.Entity)
	if tmpErr != nil {
		return nil, tmpErr
	}
	defer func() {
		if tmpErr := db.dropTable(ctx, entityDfWithFeatureName); tmpErr != nil {
			err = tmpErr
		}
	}()

	entityDfName, tmpErr := db.createAndImportTableEntityDf(ctx, opt.EntityRows, opt.Entity)
	if tmpErr != nil {
		return nil, tmpErr
	}
	defer func() {
		if tmpErr := db.dropTable(ctx, entityDfName); tmpErr != nil {
			err = tmpErr
		}
	}()

	// Step 1: iterate each table range, get result
	joinQuery := `
		INSERT INTO %s(unique_key, entity_key, unix_time, %s)
		SELECT
			CONCAT(l.entity_key, ',', l.unix_time) AS unique_key,
			l.entity_key AS entity_key,
			l.unix_time AS unix_time,
			%s
		FROM %s AS l
		LEFT JOIN %s AS r
		ON l.entity_key = r.%s
		WHERE l.unix_time >= $1 AND l.unix_time < $2;
	`
	featureNamesStr := buildFeatureNameStr(opt.Features)
	for _, r := range opt.RevisionRanges {
		_, tmpErr := db.ExecContext(ctx, fmt.Sprintf(joinQuery, entityDfWithFeatureName, featureNamesStr, featureNamesStr, entityDfName, r.DataTable, opt.Entity.Name), r.MinRevision, r.MaxRevision)
		if tmpErr != nil {
			return nil, tmpErr
		}
	}

	// Step 2: get rows from entity_df_with_features table
	resultQuery := fmt.Sprintf(`SELECT * FROM %s`, entityDfWithFeatureName)
	rows, tmpErr := db.QueryxContext(ctx, resultQuery)
	if tmpErr != nil {
		return nil, tmpErr
	}
	defer rows.Close()

	dataMap, err = getFeatureValueMapFromRows(rows, "unique_key")
	return dataMap, err
}

func buildFeatureNameStr(features []*types.RichFeature) string {
	featureNames := make([]string, 0, len(features))
	for _, f := range features {
		featureNames = append(featureNames, f.Name)
	}
	return strings.Join(featureNames, ", ")
}

func getFeatureValueMapFromRows(rows *sqlx.Rows, entityName string) (map[string]dbutil.RowMap, error) {
	featureValueMap := make(map[string]dbutil.RowMap)
	for rows.Next() {
		rowMap := make(dbutil.RowMap)
		if err := rows.MapScan(rowMap); err != nil {
			return nil, err
		}
		entityKey, ok := rowMap[entityName]
		if !ok {
			return nil, fmt.Errorf("missing column %s", entityName)
		}
		delete(rowMap, entityName)
		featureValueMap[cast.ToString(entityKey)] = rowMap
	}
	return featureValueMap, nil
}
