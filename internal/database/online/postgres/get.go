package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	featureNames := opt.FeatureList.Names()
	tableName := getOnlineBatchTableName(opt.RevisionID)
	query := fmt.Sprintf(`SELECT %s FROM "%s" WHERE "%s" = $1`, dbutil.Quote(`"`, featureNames...), tableName, opt.Entity.Name)

	record, err := db.QueryRowxContext(ctx, query, opt.EntityKey).SliceScan()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UndefinedTable {
				return nil, nil
			}
		}
		return nil, err
	}

	rs, err := deserializeIntoRowMap(record, opt.FeatureList)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	featureNames := opt.FeatureList.Names()
	tableName := getOnlineBatchTableName(opt.RevisionID)
	query := fmt.Sprintf(`SELECT "%s", %s FROM "%s" WHERE "%s" in (?);`, opt.Entity.Name, dbutil.Quote(`"`, featureNames...), tableName, opt.Entity.Name)
	sql, args, err := sqlx.In(query, opt.EntityKeys)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryxContext(ctx, db.Rebind(sql), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return getFeatureValueMapFromRows(rows, opt.FeatureList)
}

func getFeatureValueMapFromRows(rows *sqlx.Rows, features types.FeatureList) (map[string]dbutil.RowMap, error) {
	featureValueMap := make(map[string]dbutil.RowMap)
	for rows.Next() {
		record, err := rows.SliceScan()
		if err != nil {
			return nil, err
		}
		entityKey, values := record[0].(string), record[1:]
		rowMap, err := deserializeIntoRowMap(values, features)
		if err != nil {
			return nil, err
		}
		featureValueMap[entityKey] = rowMap
	}
	return featureValueMap, nil
}

func deserializeIntoRowMap(values []interface{}, features types.FeatureList) (dbutil.RowMap, error) {
	rs := map[string]interface{}{}
	for i := range values {
		rs[features[i].Name] = values[i]
	}
	return rs, nil
}
