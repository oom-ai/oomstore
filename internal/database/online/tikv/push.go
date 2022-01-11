package tikv

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/kvutil"
	"github.com/pkg/errors"
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

		serializedFeatureID, err := kvutil.SerializeByValue(opt.Features[i].ID)
		if err != nil {
			return err
		}

		serializedFeatureValue, err := dbutil.SerializeByValueType(value, opt.Features[i].ValueType, types.BackendTiKV)
		if err != nil {
			return err
		}

		putKeys = append(putKeys, getKeyOfStreamFeature(serializedGroupID, serializedEntityKey, serializedFeatureID))
		putVals = append(putVals, []byte(serializedFeatureValue))
	}

	// We don't expire keys using TTL
	err = db.BatchPut(ctx, putKeys, putVals, []uint64{})
	return errors.WithStack(err)
}
