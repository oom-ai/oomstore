package redis

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	var seq int64
	pipe := db.Pipeline()
	defer pipe.Close()

	for record := range opt.ExportStream {
		if record.Error != nil {
			return record.Error
		}
		if len(record.Record) != len(opt.Features)+1 {
			return fmt.Errorf("field count not matched, expected %d, got %d", len(opt.Features)+1, len(record.Record))
		}

		entityKey, values := record.Record[0], record.Record[1:]
		var key string
		var err error
		if opt.Group.Category == types.CategoryBatch {
			key, err = serializeRedisKeyForBatchFeature(*opt.RevisionID, entityKey)
		} else {
			key, err = serializeRedisKeyForStreamFeature(opt.Group.ID, entityKey)
		}
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
			if _, err = pipe.Exec(ctx); err != nil {
				return errdefs.WithStack(err)
			}
		}
	}

	if seq%PipelineBatchSize != 0 {
		if _, err := pipe.Exec(ctx); err != nil {
			return errdefs.WithStack(err)
		}
	}
	return nil
}
