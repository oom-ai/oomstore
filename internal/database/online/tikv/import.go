package tikv

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/redis"
)

const BatchSize = 100

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	var seq int64

	serializedRevisionID, err := redis.SerializeByValue(opt.Revision.ID)
	if err != nil {
		return err
	}

	var serializedFeatureIDs []string
	for _, feature := range opt.FeatureList {
		serializedFeatureID, err := redis.SerializeByValue(feature.ID)
		if err != nil {
			return err
		}
		serializedFeatureIDs = append(serializedFeatureIDs, serializedFeatureID)
	}

	// For rawkv.Client.BatchPut(ctx, putKeys, values)
	var putKeys [][]byte
	var putVals [][]byte

	for record := range opt.ExportStream {
		if len(record) != len(opt.FeatureList)+1 {
			return fmt.Errorf("field count not matched, expected %d, got %d", len(opt.FeatureList)+1, len(record))
		}

		entityKey, featureValues := record[0], record[1:]

		serializedEntityKey, err := redis.SerializeByValue(entityKey)
		if err != nil {
			return err
		}

		for i := range opt.FeatureList {
			// omit nil feature value
			if featureValues[i] == nil {
				continue
			}

			serializedFeatureValue, err := redis.SerializeByTag(featureValues[i], opt.FeatureList[i].ValueType)
			if err != nil {
				return err
			}

			putKeys = append(putKeys, getKey(serializedRevisionID, serializedEntityKey, serializedFeatureIDs[i]))
			putVals = append(putVals, []byte(serializedFeatureValue))
		}

		seq++
		if seq%BatchSize == 0 {
			// We don't expire keys using TTL
			if err = db.BatchPut(ctx, putKeys, putVals, []uint64{}); err != nil {
				return err
			}
			// Reset the slices
			putKeys, putVals = nil, nil
		}
	}

	if seq%BatchSize != 0 {
		// We don't expire keys using TTL
		if err := db.BatchPut(ctx, putKeys, putVals, []uint64{}); err != nil {
			return err
		}
	}

	if opt.ExportError != nil {
		return <-opt.ExportError
	}

	return nil
}
