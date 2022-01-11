package cassandra

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	var tableName string
	if opt.Group.Category == types.CategoryBatch {
		tableName = sqlutil.OnlineBatchTableName(*opt.RevisionID)
	} else {
		tableName = sqlutil.OnlineStreamTableName(opt.Group.ID)
	}

	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ?`,
		strings.Join(opt.Features.Names(), ","),
		tableName,
		opt.Entity.Name,
	)

	scan := make(map[string]interface{})
	if err := db.Query(query, opt.EntityKey).WithContext(ctx).MapScan(scan); err != nil {
		if err == gocql.ErrNotFound || isTableNotFoundError(err, tableName) {
			return scan, nil
		}
		return nil, errors.WithStack(err)
	}

	rs := make(map[string]interface{}, len(scan))
	for _, feature := range opt.Features {
		deserializedValue, _ := dbutil.DeserializeByValueType(scan[feature.Name], feature.ValueType, types.BackendCassandra)
		rs[feature.FullName] = deserializedValue
	}
	return rs, nil
}

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	var (
		tableName    string
		placeholders = getPlaceholders(len(opt.EntityKeys))
	)

	if opt.Group.Category == types.CategoryBatch {
		tableName = sqlutil.OnlineBatchTableName(*opt.RevisionID)
	} else {
		tableName = sqlutil.OnlineStreamTableName(opt.Group.ID)
	}

	query := fmt.Sprintf(`SELECT %s, %s FROM %s WHERE %s in (%s)`,
		opt.Entity.Name,
		strings.Join(opt.Features.Names(), ","),
		tableName,
		opt.Entity.Name,
		placeholders,
	)

	rs := make(map[string]dbutil.RowMap)
	scan, err := db.Query(query, toInterfaceSlice(opt.EntityKeys)...).
		WithContext(ctx).
		Iter().SliceMap()
	if err != nil {
		if err == gocql.ErrNotFound || isTableNotFoundError(err, tableName) {
			return rs, nil
		}
		return nil, errors.WithStack(err)
	}

	for _, s := range scan {
		entityKey, value := deserializeIntoRowMap(s, opt.Entity.Name, opt.Features)
		rs[entityKey] = value
	}
	return rs, nil
}

func deserializeIntoRowMap(values map[string]interface{}, entityName string, features types.FeatureList) (string, dbutil.RowMap) {
	entityKey := values[entityName].(string)

	rs := make(dbutil.RowMap)
	for _, feature := range features {
		deserializedValue, _ := dbutil.DeserializeByValueType(values[feature.Name], feature.ValueType, types.BackendCassandra)
		rs[feature.FullName] = deserializedValue
	}
	return entityKey, rs
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
