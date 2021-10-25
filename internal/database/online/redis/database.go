package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/internal/database/online"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

var _ online.Store = &DB{}

const PipelineBatchSize = 10
const SeralizeIntBase = 36

type DB struct {
	*redis.Client
}

type RedisOpt struct {
	Host string
	Port int
	Pass string
	DB   int
}

func Open(opt *RedisOpt) *DB {
	redisOpt := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		Password: opt.Pass,
		DB:       opt.DB,
	}
	return &DB{redis.NewClient(&redisOpt)}
}

func (s *DB) GetFeatureValues(ctx context.Context, dataTable, entityName, entityKey string, featureNames []string) (database.RowMap, error) {
	panic("unimplemented")
}

func (s *DB) GetFeatureValuesWithMultiEntityKeys(ctx context.Context, dataTable, entityName string, entityKeys, featureNames []string) (map[string]database.RowMap, error) {
	panic("unimplemented")
}

func (s *DB) SinkFeatureValuesStream(ctx context.Context, stream <-chan []interface{}, features []*types.Feature, revision *types.Revision) error {
	var seq int64
	pipe := s.Pipeline()
	defer pipe.Close()

	for row := range stream {
		if len(row) != len(features)+1 {
			return fmt.Errorf("field count not matched, expected %d, got %d", len(features)+1, len(row))
		}

		revisionId, err := Seralize(revision.ID)
		if err != nil {
			return err
		}

		entityKey, err := Seralize(row[0])
		if err != nil {
			return err
		}

		values := row[1:]
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
		return strconv.FormatBool(s), nil

	default:
		return "", fmt.Errorf("unable to seralize %#v of type %T to string", i, i)
	}
}
