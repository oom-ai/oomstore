package tikv

import (
	"context"

	"github.com/pkg/errors"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/kvutil"
)

const BatchSize = 100

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	var seq int64

	serializedRevisionID, err := kvutil.SerializeByValue(opt.Revision.ID)
	if err != nil {
		return err
	}

	var serializedFeatureIDs []string
	for _, feature := range opt.Features {
		serializedFeatureID, err := kvutil.SerializeByValue(feature.ID)
		if err != nil {
			return err
		}
		serializedFeatureIDs = append(serializedFeatureIDs, serializedFeatureID)
	}

	// For rawkv.Client.BatchPut(ctx, putKeys, values)
	var putKeys [][]byte
	var putVals [][]byte

	for record := range opt.ExportStream {
		if len(record) != len(opt.Features)+1 {
			return errors.Errorf("field count not matched, expected %d, got %d", len(opt.Features)+1, len(record))
		}

		entityKey, featureValues := record[0], record[1:]

		serializedEntityKey, err := kvutil.SerializeByValue(entityKey)
		if err != nil {
			return err
		}

		for i := range opt.Features {
			// omit nil feature value
			if featureValues[i] == nil {
				continue
			}

			serializedFeatureValue, err := kvutil.SerializeByValueType(featureValues[i], opt.Features[i].ValueType)
			if err != nil {
				return err
			}

			putKeys = append(putKeys, getKeyOfBatchFeature(serializedRevisionID, serializedEntityKey, serializedFeatureIDs[i]))
			putVals = append(putVals, []byte(serializedFeatureValue))
		}

		seq++
		if seq%BatchSize == 0 {
			// We don't expire keys using TTL
			if err = db.BatchPut(ctx, putKeys, putVals, []uint64{}); err != nil {
				return errors.WithStack(err)
			}
			// Reset the slices
			putKeys, putVals = nil, nil
		}
	}

	if seq%BatchSize != 0 {
		// We don't expire keys using TTL
		if err := db.BatchPut(ctx, putKeys, putVals, []uint64{}); err != nil {
			return errors.WithStack(err)
		}
	}

	if opt.ExportError != nil {
		return <-opt.ExportError
	}

	return nil
}
