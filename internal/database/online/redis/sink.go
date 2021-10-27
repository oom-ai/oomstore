package redis

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/onestore/types"
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	var seq int64
	pipe := db.Pipeline()
	defer pipe.Close()

	for item := range opt.Stream {
		if item.Error != nil {
			return item.Error
		}
		record := item.Record
		if len(record) != len(opt.Features)+1 {
			return fmt.Errorf("field count not matched, expected %d, got %d", len(opt.Features)+1, len(record))
		}

		entityKey, values := record[0], record[1:]

		key, err := SerializeRedisKey(opt.Revision.ID, entityKey)
		if err != nil {
			return err
		}

		featureValues := make(map[string]string)
		for i := range opt.Features {
			// omit nil feature value
			if values[i] == nil {
				continue
			}
			featureValue, err := SerializeByTag(values[i], opt.Features[i].ValueType)
			if err != nil {
				return err
			}

			featureId, err := SerializeByValue(opt.Features[i].ID)
			if err != nil {
				return err
			}
			featureValues[featureId] = featureValue
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
	return nil
}

func (db *DB) Purge(ctx context.Context, revision *types.Revision) error {
	prefix, err := SerializeByValue(revision.ID)
	if err != nil {
		return nil
	}
	pattern := prefix + ":*"

	var cursor uint64
	for {
		keys, cursor, err := db.Scan(ctx, cursor, pattern, PipelineBatchSize).Result()
		if err != nil {
			return err
		}

		if _, err = db.Del(ctx, keys...).Result(); err != nil {
			return err
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}
