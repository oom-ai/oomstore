package redis

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func TestImport(t *testing.T) {
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
		stream <- &types.RawFeatureValueRecord{
			Record: []interface{}{"3215", int16(18), "F"},
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

	{
		field, err := SerializeByValue(feature1.ID)
		require.NoError(t, err)
		value, err := SerializeByValue(int16(18))
		require.NoError(t, err)
		key, err := SerializeRedisKey(revision.ID, "3215")
		require.NoError(t, err)
		age, err := store.HGet(ctx, key, field).Result()
		require.NoError(t, err)
		require.Equal(t, age, value)
	}

	{
		field, err := SerializeByValue(feature2.ID)
		require.NoError(t, err)
		value, err := SerializeByValue("F")
		require.NoError(t, err)
		key, err := SerializeRedisKey(revision.ID, "3215")
		require.NoError(t, err)
		gender, err := store.HGet(ctx, key, field).Result()
		require.NoError(t, err)
		require.Equal(t, gender, value)
	}
}
