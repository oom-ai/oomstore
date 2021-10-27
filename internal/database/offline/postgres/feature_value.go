package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"

	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) GetPointInTimeFeatureValues(ctx context.Context, entity *types.Entity,
	revisionRanges []*types.RevisionRange, features []*types.RichFeature, entityRows []types.EntityRow) (dataMap map[string]database.RowMap, err error) {
	if len(features) == 0 {
		return make(map[string]database.RowMap), nil
	}

	// Step 0: prepare temporary tables
	entityDfWithFeatureName, tmpErr := db.createTableEntityDfWithFeatures(ctx, features, entity)
	if tmpErr != nil {
		return nil, tmpErr
	}
	defer func() {
		if tmpErr := db.dropTable(ctx, entityDfWithFeatureName); tmpErr != nil {
			err = tmpErr
		}
	}()

	entityDfName, tmpErr := db.createAndImportTableEntityDf(ctx, entityRows, entity)
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
	featureNamesStr := buildFeatureNameStr(features)
	for _, r := range revisionRanges {
		_, tmpErr := db.ExecContext(ctx, fmt.Sprintf(joinQuery, entityDfWithFeatureName, featureNamesStr, featureNamesStr, entityDfName, r.DataTable, entity.Name), r.MinRevision, r.MaxRevision)
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

func getFeatureValueMapFromRows(rows *sqlx.Rows, entityName string) (map[string]database.RowMap, error) {
	featureValueMap := make(map[string]database.RowMap)
	for rows.Next() {
		rowMap := make(database.RowMap)
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
