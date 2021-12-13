package cassandra

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocql/gocql"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/internal/database/online/sqlutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	tableName := sqlutil.OnlineTableName(opt.RevisionID)

	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ?`,
		strings.Join(opt.FeatureList.Names(), ","),
		tableName,
		opt.Entity.Name,
	)

	rs := make(map[string]interface{})
	if err := db.Query(query, opt.EntityKey).WithContext(ctx).MapScan(rs); err != nil {
		if err == gocql.ErrNotFound || isTableNotFoundError(err, tableName) {
			return nil, nil
		}
		return nil, err
	}

	for k, v := range rs {
		rs[k] = deserializeString(v)
	}
	return rs, nil
}

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	var (
		tableName    = sqlutil.OnlineTableName(opt.RevisionID)
		placeholders = getPlaceholders(len(opt.EntityKeys))
	)

	query := fmt.Sprintf(`SELECT %s, %s FROM %s WHERE %s in (%s)`,
		opt.Entity.Name,
		strings.Join(opt.FeatureList.Names(), ","),
		tableName,
		opt.Entity.Name,
		placeholders,
	)

	rs := make(map[string]dbutil.RowMap)
	slice, err := db.Query(query, toInterfaceSlice(opt.EntityKeys)...).
		WithContext(ctx).
		Iter().SliceMap()
	if err != nil {
		if err == gocql.ErrNotFound || isTableNotFoundError(err, tableName) {
			return nil, nil
		}
		return nil, err
	}

	for _, s := range slice {
		entityKey, value := deserializeIntoRowMap(s, opt.Entity.Name, opt.FeatureList)
		rs[entityKey] = value

	}
	return rs, nil
}

func deserializeString(i interface{}) interface{} {
	switch i.(type) {
	case string:
		if i == "" {
			return nil
		}
	}
	return i
}

func deserializeIntoRowMap(values map[string]interface{}, entityName string, features types.FeatureList) (string, dbutil.RowMap) {
	entityKey := values[entityName].(string)
	delete(values, entityName)

	for k, v := range values {
		values[k] = deserializeString(v)
	}
	return entityKey, values
}

func isTableNotFoundError(err error, tableName string) bool {
	if err == nil {
		return false
	}
	return err.Error() == fmt.Sprintf("table %s does not exist", tableName)
}

func getPlaceholders(length int) string {
	p := make([]string, length)
	for i := 0; i < length; i++ {
		p[i] = "?"
	}
	return strings.Join(p, ",")
}

func toInterfaceSlice(s []string) []interface{} {
	rs := make([]interface{}, 0, len(s))
	for _, key := range s {
		rs = append(rs, key)
	}
	return rs
}
