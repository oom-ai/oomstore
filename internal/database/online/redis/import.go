package redis

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/kvutil"
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	var seq int64
	pipe := db.Pipeline()
	defer pipe.Close()

	for record := range opt.ExportStream {
		if len(record) != len(opt.FeatureList)+1 {
			return fmt.Errorf("field count not matched, expected %d, got %d", len(opt.FeatureList)+1, len(record))
		}

		entityKey, values := record[0], record[1:]
		key, err := serializeRedisKeyForBatchFeature(opt.Revision.ID, entityKey)
		if err != nil {
			return err
		}

		featureValues := make(map[string]string)
		for i := range opt.FeatureList {
			// omit nil feature value
			if values[i] == nil {
				continue
			}
			featureValue, err := kvutil.SerializeByValueType(values[i], opt.FeatureList[i].ValueType)
			if err != nil {
				return err
			}

			featureID, err := kvutil.SerializeByValue(opt.FeatureList[i].ID)
			if err != nil {
				return err
			}
			featureValues[featureID] = featureValue
		}

		pipe.HSet(ctx, key, featureValues)

		seq++
		if seq%PipelineBatchSize == 0 {
			if _, err := pipe.Exec(ctx); err != nil {
				return err
			}
		}
	}

	if seq%PipelineBatchSize != 0 {
		if _, err := pipe.Exec(ctx); err != nil {
			return err
		}
	}
	if opt.ExportError != nil {
		return <-opt.ExportError
	}
	return nil
}
