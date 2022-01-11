package redis

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/pkg/errors"

	"github.com/oom-ai/oomstore/internal/database/online"
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	var seq int64
	pipe := db.Pipeline()
	defer pipe.Close()

	for record := range opt.ExportStream {
		if len(record) != len(opt.Features)+1 {
			return fmt.Errorf("field count not matched, expected %d, got %d", len(opt.Features)+1, len(record))
		}

		entityKey, values := record[0], record[1:]
		key, err := serializeRedisKeyForBatchFeature(opt.Revision.ID, entityKey)
		if err != nil {
			return err
		}

		featureValues := make(map[string]string)
		for i := range opt.Features {
			// omit nil feature value
			if values[i] == nil {
				continue
			}
			featureValue, err := dbutil.SerializeByValueType(values[i], opt.Features[i].ValueType, Backend)
			if err != nil {
				return err
			}

			featureID, err := dbutil.SerializeByValue(opt.Features[i].ID, Backend)
			if err != nil {
				return err
			}
			featureValues[featureID] = featureValue.(string)
		}

		pipe.HSet(ctx, key, featureValues)

		seq++
		if seq%PipelineBatchSize == 0 {
			if _, err := pipe.Exec(ctx); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	if seq%PipelineBatchSize != 0 {
		if _, err := pipe.Exec(ctx); err != nil {
			return errors.WithStack(err)
		}
	}
	if opt.ExportError != nil {
		return <-opt.ExportError
	}
	return nil
}
