package tikv

import (
	"context"

	"github.com/pkg/errors"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type input struct {
	entityKey string
	feature   types.Feature
}

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	// Proxy to MultiGet
	res, err := db.MultiGet(ctx, online.MultiGetOpt{
		Entity:     opt.Entity,
		Group:      opt.Group,
		RevisionID: opt.RevisionID,
		EntityKeys: []string{opt.EntityKey},
		Features:   opt.Features,
	})
	if err != nil {
		return nil, err
	}
	if rowMap, ok := res[opt.EntityKey]; !ok {
		return make(dbutil.RowMap), nil
	} else {
		return rowMap, nil
	}
}

func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	var (
		allBatch            bool
		serializedPrefixKey string
		err                 error
	)

	allBatch = (opt.Group.Category == types.CategoryBatch)

	if allBatch {
		serializedPrefixKey, err = dbutil.SerializeByValue(*opt.RevisionID, Backend)
	} else {
		serializedPrefixKey, err = dbutil.SerializeByValue(opt.Group.ID, Backend)
	}
	if err != nil {
		return nil, err
	}

	var serializedEntityKeys []string
	for _, entityKey := range opt.EntityKeys {
		serializedEntityKey, err := dbutil.SerializeByValue(entityKey, Backend)
		if err != nil {
			return nil, err
		}
		serializedEntityKeys = append(serializedEntityKeys, serializedEntityKey)
	}

	var serializedFeatureIDs []string
	for _, feature := range opt.Features {
		serializedFeatureID, err := dbutil.SerializeByValue(feature.ID, Backend)
		if err != nil {
			return nil, err
		}
		serializedFeatureIDs = append(serializedFeatureIDs, serializedFeatureID)
	}

	// What rawkv.Client needs
	var keys [][]byte
	for _, serializedEntityKey := range serializedEntityKeys {
		for _, serializedFeatureID := range serializedFeatureIDs {
			if allBatch {
				keys = append(keys, getKeyOfBatchFeature(serializedPrefixKey, serializedEntityKey, serializedFeatureID))
			} else {
				keys = append(keys, getKeyOfStreamFeature(serializedPrefixKey, serializedEntityKey, serializedFeatureID))
			}
		}
	}

	// Result order is the same as input order
	batchGetResult, err := db.BatchGet(ctx, keys)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// What we need to align with the result
	var inputs []input
	for _, entityKey := range opt.EntityKeys {
		for _, feature := range opt.Features {
			inputs = append(inputs, input{entityKey, *feature})
		}
	}

	res := make(map[string]dbutil.RowMap)
	for i, v := range batchGetResult {
		if v == nil {
			continue
		}
		deserializedValue, err := dbutil.DeserializeByValueType(string(v), inputs[i].feature.ValueType, Backend)
		if err != nil {
			return nil, err
		}
		entityKey := inputs[i].entityKey
		if _, ok := res[entityKey]; !ok {
			res[entityKey] = make(dbutil.RowMap)
		}
		res[entityKey][inputs[i].feature.FullName] = deserializedValue
	}
	return res, nil
}
