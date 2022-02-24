package tikv

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/online"
)

const BatchSize = 100

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	var seq int64
	var err error
	var serializedRevisionID, serializedGroupID string

	if opt.Group.Category == types.CategoryBatch {
		serializedRevisionID, err = dbutil.SerializeByValue(*opt.RevisionID, Backend)
	} else {
		serializedGroupID, err = dbutil.SerializeByValue(opt.Group.ID, Backend)
	}
	if err != nil {
		return err
	}

	var serializedFeatureIDs []string
	for _, feature := range opt.Features {
		serializedFeatureID, err2 := dbutil.SerializeByValue(feature.ID, Backend)
		if err2 != nil {
			return err2
		}
		serializedFeatureIDs = append(serializedFeatureIDs, serializedFeatureID)
	}

	// For rawkv.Client.BatchPut(ctx, putKeys, values)
	var putKeys [][]byte
	var putVals [][]byte

	for record := range opt.ExportStream {
		if record.Error != nil {
			return record.Error
		}

		if len(record.Record) != len(opt.Features)+1 {
			return errdefs.Errorf("field count not matched, expected %d, got %d", len(opt.Features)+1, len(record.Record))
		}

		entityKey, featureValues := record.Record[0], record.Record[1:]

		serializedEntityKey, err := dbutil.SerializeByValue(entityKey, Backend)
		if err != nil {
			return err
		}

		for i := range opt.Features {
			// omit nil feature value
			if featureValues[i] == nil {
				continue
			}

			serializedFeatureValue, err2 := dbutil.SerializeByValueType(featureValues[i], opt.Features[i].ValueType, Backend)
			if err2 != nil {
				return err
			}
			if opt.Group.Category == types.CategoryBatch {
				putKeys = append(putKeys, getKeyOfBatchFeature(serializedRevisionID, serializedEntityKey, serializedFeatureIDs[i]))
			} else {
				putKeys = append(putKeys, getKeyOfStreamFeature(serializedGroupID, serializedEntityKey, serializedFeatureIDs[i]))
			}
			putVals = append(putVals, []byte(serializedFeatureValue.(string)))
		}

		seq++
		if seq%BatchSize == 0 {
			// We don't expire keys using TTL
			if err = db.BatchPut(ctx, putKeys, putVals, []uint64{}); err != nil {
				return errdefs.WithStack(err)
			}
			// Reset the slices
			putKeys, putVals = nil, nil
		}
	}

	if seq%BatchSize != 0 {
		// We don't expire keys using TTL
		if err := db.BatchPut(ctx, putKeys, putVals, []uint64{}); err != nil {
			return errdefs.WithStack(err)
		}
	}
	return nil
}
