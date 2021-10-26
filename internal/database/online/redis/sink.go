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

		revisionId, err := Seralize(revision.ID)
		if err != nil {
			return err
		}

		entityKey, err := Seralize(record[0])
		if err != nil {
			return err
		}

		values := record[1:]
		featureValues := make(map[string]string)
		for i := range features {
			// omit nil feature value
			if values[i] == nil {
				continue
			}
			featureValue, err := Seralize(values[i])
			if err != nil {
				return err
			}

			featureId, err := Seralize(features[i].ID)
			if err != nil {
				return err
			}
			featureValues[featureId] = featureValue
		}
		pipe.HSet(ctx, fmt.Sprintf("%s:%s", revisionId, entityKey), featureValues)

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

func (db *DB) DeprecateFeatureValues(ctx context.Context, tableName string) error {
	panic("implement me")
}
