package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) SinkFeatureValuesStream(ctx context.Context, stream <-chan *types.RawFeatureValueRecord, features []*types.Feature, revision *types.Revision) error {
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

// seralize feature values into compact string
func Seralize(i interface{}) (string, error) {
	switch s := i.(type) {
	case string:
		return s, nil
	case []byte:
		return string(s), nil

	case int:
		return strconv.FormatInt(int64(s), SeralizeIntBase), nil
	case int64:
		return strconv.FormatInt(int64(s), SeralizeIntBase), nil
	case int32:
		return strconv.FormatInt(int64(s), SeralizeIntBase), nil
	case int16:
		return strconv.FormatInt(int64(s), SeralizeIntBase), nil
	case int8:
		return strconv.FormatInt(int64(s), SeralizeIntBase), nil

	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil

	case uint:
		return strconv.FormatUint(uint64(s), SeralizeIntBase), nil
	case uint64:
		return strconv.FormatUint(uint64(s), SeralizeIntBase), nil
	case uint32:
		return strconv.FormatUint(uint64(s), SeralizeIntBase), nil
	case uint16:
		return strconv.FormatUint(uint64(s), SeralizeIntBase), nil
	case uint8:
		return strconv.FormatUint(uint64(s), SeralizeIntBase), nil

	case time.Time:
		return Seralize(s.UnixMilli())
	case bool:
		if s {
			return "1", nil
		} else {
			return "0", nil
		}

	default:
		return "", fmt.Errorf("unable to seralize %#v of type %T to string", i, i)
	}
}
