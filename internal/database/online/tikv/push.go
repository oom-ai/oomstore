package tikv

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/kvutil"
)

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	serializedEntityKey, err := kvutil.SerializeByValue(opt.EntityKey)
	if err != nil {
		return err
	}
	serializedGroupID, err := kvutil.SerializeByValue(opt.GroupID)
	if err != nil {
		return err
	}

	// For rawkv.Client.BatchPut(ctx, putKeys, values)
	var putKeys [][]byte
	var putVals [][]byte

	for i, value := range opt.FeatureValues {
		// omit nil feature value
		if value == nil {
			continue
		}

		serializedFeatureID, err := kvutil.SerializeByValue(opt.FeatureList[i].ID)
		if err != nil {
			return err
		}

		serializedFeatureValue, err := kvutil.SerializeByValueType(value, opt.FeatureList[i].ValueType)
		if err != nil {
			return err
		}

		putKeys = append(putKeys, getKeyOfStreamFeature(serializedGroupID, serializedEntityKey, serializedFeatureID))
		putVals = append(putVals, []byte(serializedFeatureValue))
	}

	// We don't expire keys using TTL
	return db.BatchPut(ctx, putKeys, putVals, []uint64{})
}
