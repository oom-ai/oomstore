package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/spf13/cast"
)

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	featureNames := []string{}
	for _, f := range opt.Features {
		featureNames = append(featureNames, f.Name)
	}
	query := fmt.Sprintf(`SELECT "%s",%s FROM %s WHERE "%s" = $1`, opt.EntityName, strings.Join(featureNames, ","), opt.DataTable, opt.EntityName)
	rs := make(dbutil.RowMap)

	if err := db.QueryRowxContext(ctx, query, opt.EntityKey).MapScan(rs); err != nil {
		return nil, err
	}
	return rs, nil
}

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	featureNames := []string{}
	for _, f := range opt.Features {
		featureNames = append(featureNames, f.Name)
	}
	query := fmt.Sprintf(`SELECT "%s", %s FROM %s WHERE "%s" in (?);`, opt.EntityName, strings.Join(featureNames, ","), opt.DataTable, opt.EntityName)
	sql, args, err := sqlx.In(query, opt.EntityKeys)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryxContext(ctx, db.Rebind(sql), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return getFeatureValueMapFromRows(rows, opt.EntityName)
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
