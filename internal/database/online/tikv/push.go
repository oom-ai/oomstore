package tikv

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/errdefs"
)

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	serializedEntityKey, err := dbutil.SerializeByValue(opt.EntityKey, Backend)
	if err != nil {
		return err
	}
	serializedGroupID, err := dbutil.SerializeByValue(opt.GroupID, Backend)
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

		serializedFeatureID, err := dbutil.SerializeByValue(opt.Features[i].ID, Backend)
		if err != nil {
			return err
		}

		serializedFeatureValue, err := dbutil.SerializeByValueType(value, opt.Features[i].ValueType, Backend)
		if err != nil {
			return err
		}

		putKeys = append(putKeys, getKeyOfStreamFeature(serializedGroupID, serializedEntityKey, serializedFeatureID))
		putVals = append(putVals, []byte(serializedFeatureValue.(string)))
	}

	// We don't expire keys using TTL
	err = db.BatchPut(ctx, putKeys, putVals, []uint64{})
	return errdefs.WithStack(err)
}
