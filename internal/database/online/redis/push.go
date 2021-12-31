package redis

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/kvutil"
)

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	key, err := serializeRedisKeyForStreamFeature(opt.GroupID, opt.EntityKey)
	if err != nil {
		return err
	}

	featureValues := make(map[string]string)
	for i, value := range opt.FeatureValues {
		// omit nil feature value
		if value == nil {
			continue
		}

		featureValue, err := kvutil.SerializeByValueType(value, opt.FeatureList[i].ValueType)
		if err != nil {
			return err
		}

		featureID, err := kvutil.SerializeByValue(opt.FeatureList[i].ID)
		if err != nil {
			return err
		}
		featureValues[featureID] = featureValue
	}

	db.HSet(ctx, key, featureValues)
	return nil
}
