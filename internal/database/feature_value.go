package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
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
		featureValueMap[entityKey.(string)] = rowMap
	}
	return featureValueMap, nil
}
