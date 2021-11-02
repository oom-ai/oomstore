package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_redis"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func prepare() (context.Context, *DB) {
	ctx := context.Background()
	store := Open(&runtime_redis.RedisDbOpt)
	if _, err := store.FlushDB(ctx).Result(); err != nil {
		panic(err)
	}
	return ctx, store
}

func importSample(t *testing.T) online.ImportOpt {
	ctx, store := prepare()

	feature1 := types.Feature{
		ID:        0,
		Name:      "age",
		GroupName: "user",
		ValueType: "int16",
	}
	feature2 := types.Feature{
		ID:        1,
		Name:      "gender",
		GroupName: "user",
		ValueType: "string",
	}
	features := types.FeatureList{&feature1, &feature2}
	revision := types.Revision{ID: 3}
	entity := types.Entity{ID: 5}
	stream := make(chan *types.RawFeatureValueRecord)
	go func() {
		defer close(stream)

		records := [][]interface{}{
			{"3215", int16(18), "F"},
			{"3216", int16(29), nil},
			{"3217", int16(44), "M"},
		}
		for _, record := range records {
			stream <- &types.RawFeatureValueRecord{
				Record: record,
			}
		}
	}()

	opt := online.ImportOpt{
		Features: features,
		Revision: &revision,
		Entity:   &entity,
		Stream:   stream,
	}

	err := store.Import(ctx, opt)
	require.NoError(t, err)

	return opt
}

func TestOpen(t *testing.T) {
	ctx, store := prepare()
	res, err := store.Ping(ctx).Result()
	assert.Nil(t, err)
	assert.Equal(t, res, "PONG")
}
