package sqlutil

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func Get(ctx context.Context, db *sqlx.DB, opt online.GetOpt, backend types.BackendType) (dbutil.RowMap, error) {
	var tableName string
	if opt.Group.Category == types.CategoryBatch {
		tableName = OnlineBatchTableName(*opt.RevisionID)
	} else {
		tableName = OnlineStreamTableName(opt.Group.ID)
	}

	featureNames := opt.Features.Names()
	qt := dbutil.QuoteFn(backend)
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ?`, qt(featureNames...), qt(tableName), qt(opt.Entity.Name))

	record, err := db.QueryRowxContext(ctx, db.Rebind(query), opt.EntityKey).SliceScan()
	if err != nil {
		if err == sql.ErrNoRows || dbutil.IsTableNotFoundError(err, backend) {
			return make(dbutil.RowMap), nil
		}
		return nil, errdefs.WithStack(err)
	}

	rs, err := deserializeIntoRowMap(record, opt.Features, backend)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// response: map[entity_key]map[feature_name]feature_value
func MultiGet(ctx context.Context, db *sqlx.DB, opt online.MultiGetOpt, backend types.BackendType) (map[string]dbutil.RowMap, error) {
	var tableName string
	if opt.Group.Category == types.CategoryBatch {
		tableName = OnlineBatchTableName(*opt.RevisionID)
	} else {
		tableName = OnlineStreamTableName(opt.Group.ID)
	}

	featureNames := opt.Features.Names()
	qt := dbutil.QuoteFn(backend)
	query := fmt.Sprintf(`SELECT %s, %s FROM %s WHERE %s in (?);`, qt(opt.Entity.Name), qt(featureNames...), qt(tableName), qt(opt.Entity.Name))
	sql, args, err := sqlx.In(query, opt.EntityKeys)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}

	rows, err := db.QueryxContext(ctx, db.Rebind(sql), args...)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	defer rows.Close()

	return getFeatureValueMapFromRows(rows, opt.Features, backend)
}

func getFeatureValueMapFromRows(rows *sqlx.Rows, features types.FeatureList, backend types.BackendType) (map[string]dbutil.RowMap, error) {
	featureValueMap := make(map[string]dbutil.RowMap)
	for rows.Next() {
		record, err := rows.SliceScan()
		if err != nil {
			return nil, errdefs.WithStack(err)
		}
		entityKey, err := dbutil.DeserializeByValueType(record[0], types.String, backend)
		if err != nil {
			return nil, errdefs.WithStack(err)
		}
		values := record[1:]
		rowMap, err := deserializeIntoRowMap(values, features, backend)
		if err != nil {
			return nil, err
		}
		featureValueMap[entityKey.(string)] = rowMap
	}
	return featureValueMap, nil
}

func deserializeIntoRowMap(values []interface{}, features types.FeatureList, backend types.BackendType) (dbutil.RowMap, error) {
	rs := map[string]interface{}{}
	for i, v := range values {
		deserializedValue, err := dbutil.DeserializeByValueType(v, features[i].ValueType, backend)
		if err != nil {
			return nil, err
		}
		rs[features[i].FullName] = deserializedValue
	}
	return rs, nil
}
