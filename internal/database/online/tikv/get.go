package tikv

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/kvutil"
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
		RevisionID: opt.RevisionID,
		EntityKeys: []string{opt.EntityKey},
		Features:   opt.Features,
	})
	if err != nil {
		return nil, err
	}
	if rowMap, ok := res[opt.EntityKey]; !ok {
		return nil, nil
	} else {
		return rowMap, nil
	}
}

func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	serializedRevisionID, err := kvutil.SerializeByValue(opt.RevisionID)
	if err != nil {
		return nil, err
	}

	var serializedEntityKeys []string
	for _, entityKey := range opt.EntityKeys {
		serializedEntityKey, err := kvutil.SerializeByValue(entityKey)
		if err != nil {
			return nil, err
		}
		serializedEntityKeys = append(serializedEntityKeys, serializedEntityKey)
	}

	var serializedFeatureIDs []string
	for _, feature := range opt.Features {
		serializedFeatureID, err := kvutil.SerializeByValue(feature.ID)
		if err != nil {
			return nil, err
		}
		serializedFeatureIDs = append(serializedFeatureIDs, serializedFeatureID)
	}

	// What rawkv.Client needs
	var keys [][]byte
	for _, serializedEntityKey := range serializedEntityKeys {
		for _, serializedFeatureID := range serializedFeatureIDs {
			keys = append(keys, getKeyOfBatchFeature(serializedRevisionID, serializedEntityKey, serializedFeatureID))
		}
	}

	// Result order is the same as input order
	batchGetResult, err := db.BatchGet(ctx, keys)
	if err != nil {
		return nil, err
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
		typedValue, err := kvutil.DeserializeByValueType(string(v), inputs[i].feature.ValueType)
		if err != nil {
			return nil, err
		}
		entityKey := inputs[i].entityKey
		if _, ok := res[entityKey]; !ok {
			res[entityKey] = make(dbutil.RowMap)
		}
		res[entityKey][inputs[i].feature.FullName] = typedValue
	}
	return res, nil
}
