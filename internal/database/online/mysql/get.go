package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/internal/database/sqlutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	featureNames := opt.FeatureList.Names()
	tableName := sqlutil.OnlineTableName(opt.RevisionID)
	query := fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` = ?", dbutil.Quote("`", featureNames...), tableName, opt.Entity.Name)

	record, err := db.QueryRowxContext(ctx, db.Rebind(query), opt.EntityKey).SliceScan()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if e2, ok := err.(*mysql.MySQLError); ok {
			// https://dev.mysql.com/doc/mysql-errors/5.7/en/server-error-reference.html#error_er_no_such_table
			if e2.Number == 1146 {
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
	tableName := sqlutil.OnlineTableName(opt.RevisionID)
	query := fmt.Sprintf("SELECT `%s`, %s FROM `%s` WHERE `%s` in (?);", opt.Entity.Name, dbutil.Quote("`", featureNames...), tableName, opt.Entity.Name)
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
		entityKey, values := string(record[0].([]byte)), record[1:]
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
		value := values[i]
		if value != nil && features[i].ValueType == types.STRING {
			value = string(value.([]byte))
		}
		rs[features[i].Name] = value
	}
	return rs, nil
}
