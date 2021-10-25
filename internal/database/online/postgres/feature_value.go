package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/internal/database"
	"github.com/spf13/cast"
)

func (db *DB) GetFeatureValues(ctx context.Context, dataTable, entityName, entityKey string, revisionId int32, featureNames []string) (database.RowMap, error) {
	query := fmt.Sprintf(`SELECT "%s",%s FROM %s WHERE "%s" = $1`, entityName, strings.Join(featureNames, ","), dataTable, entityName)
	rs := make(database.RowMap)

	if err := db.QueryRowxContext(ctx, query, entityKey).MapScan(rs); err != nil {
		return nil, err
	}
	return rs, nil
}

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) GetFeatureValuesWithMultiEntityKeys(ctx context.Context, dataTable, entityName string, revisionId int32, entityKeys, featureNames []string) (map[string]database.RowMap, error) {
	query := fmt.Sprintf(`SELECT "%s", %s FROM %s WHERE "%s" in (?);`, entityName, strings.Join(featureNames, ","), dataTable, entityName)
	sql, args, err := sqlx.In(query, entityKeys)
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
