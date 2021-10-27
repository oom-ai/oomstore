package redis

import (
	"context"
	"fmt"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) SinkFeatureValuesStream(ctx context.Context, stream <-chan *types.RawFeatureValueRecord, features []*types.Feature, revision *types.Revision, entity *types.Entity) error {
	var seq int64
	pipe := db.Pipeline()
	defer pipe.Close()

	for item := range stream {
		if item.Error != nil {
			return item.Error
		}
		record := item.Record
		if len(record) != len(features)+1 {
			return fmt.Errorf("field count not matched, expected %d, got %d", len(features)+1, len(record))
		}

		entityKey, values := record[0], record[1:]

		key, err := SerializeRedisKey(revision.ID, entityKey)
		if err != nil {
			return err
		}

		featureValues := make(map[string]string)
		for i := range features {
			// omit nil feature value
			if values[i] == nil {
				continue
			}
			featureValue, err := SerializeByTag(values[i], features[i].ValueType)
			if err != nil {
				return err
			}

			featureId, err := SerializeByValue(features[i].ID)
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

func (db *DB) PurgeRevision(ctx context.Context, revision *types.Revision) error {
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
