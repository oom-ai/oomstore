package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
	"github.com/spf13/cast"
)

type RowMap = map[string]interface{}

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) GetFeatureValues(ctx context.Context, dataTable, entityName string, entityKeys, featureNames []string) (map[string]RowMap, error) {
	marks := []string{}
	for range featureNames {
		marks = append(marks, "?")
	}

	query := fmt.Sprintf("SELECT ?, %s FROM %s WHERE ? in (?);", strings.Join(marks, ","), dataTable)
	sql, args, err := sqlx.In(query, entityName, featureNames, entityKeys)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryxContext(ctx, db.Rebind(sql), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return getFeatureValueMapFromRows(rows, entityName)
}

func getFeatureValueMapFromRows(rows *sqlx.Rows, entityName string) (map[string]RowMap, error) {
	featureValueMap := make(map[string]RowMap)
	for rows.Next() {
		rowMap := make(RowMap)
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

func (db *DB) GetPointInTimeFeatureValues(ctx context.Context, features []*types.RichFeature, entityRows []types.EntityRow) (dataMap map[string]RowMap, err error) {
	if len(features) == 0 {
		return make(map[string]RowMap), nil
	}
	groupName := features[0].GroupName
	entityName := features[0].EntityName

	// Step 0: prepare temporary tables
	entityDfWithFeatureName, tmpErr := db.createTableEntityDfWithFeatures(ctx, features, entityName)
	if tmpErr != nil {
		return nil, tmpErr
	}
	defer func() {
		if tmpErr := db.dropTable(ctx, entityDfWithFeatureName); tmpErr != nil {
			err = tmpErr
		}
	}()

	entityDfName, tmpErr := db.createAndImportTableEntityDf(ctx, entityRows, entityName)
	if tmpErr != nil {
		return nil, tmpErr
	}
	defer func() {
		if tmpErr := db.dropTable(ctx, entityDfName); tmpErr != nil {
			err = tmpErr
		}
	}()

	// Step 1: get table ranges
	rangeQuery := `
		SELECT
			revision AS min_revision,
			LEAD(revision, 1, ~0 >> 1) OVER w AS max_revision,
			data_table
		FROM feature_group_revision
		WHERE group_name = ?
		WINDOW w AS (ORDER BY revision);
	`

	var ranges []struct {
		MinRevision int64  `db:"min_revision"`
		MaxRevision int64  `db:"max_revision"`
		DataTable   string `db:"data_table"`
	}
	if tmpErr := db.SelectContext(ctx, &ranges, rangeQuery, groupName); tmpErr != nil {
		return nil, tmpErr
	}

	// Step 2: iterate each table range, get result
	joinQuery := `
		INSERT INTO %s(unique_key, l.entity_key, l.unix_time, %s)
		SELECT
			CONCAT(l.entity_key, ",", l.unix_time) AS unique_key,
			l.entity_key, l.unix_time,
			%s
		FROM %s AS l
		LEFT JOIN %s AS r
		ON l.entity_key = r.%s
		WHERE l.unix_time >= ? AND l.unix_time < ?;
	`
	featureNamesStr := buildFeatureNameStr(features)

	for _, r := range ranges {
		_, tmpErr := db.ExecContext(ctx, fmt.Sprintf(joinQuery, entityDfWithFeatureName, featureNamesStr, featureNamesStr, entityDfName, r.DataTable, entityName), r.MinRevision, r.MaxRevision)
		if tmpErr != nil {
			return nil, tmpErr
		}
	}

	// Step 3: get rows from entity_df_with_features table
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
	return strings.Join(featureNames, " ,")
}
