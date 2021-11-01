package redis

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"gotest.tools/v3/assert"
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
	assert.NilError(t, err)

	{
		field, err := SerializeByValue(feature1.ID)
		assert.NilError(t, err)
		value, err := SerializeByValue(int16(18))
		assert.NilError(t, err)
		key, err := SerializeRedisKey(revision.ID, "3215")
		assert.NilError(t, err)
		age, err := store.HGet(ctx, key, field).Result()
		assert.NilError(t, err)
		assert.Equal(t, age, value)
	}

	{
		field, err := SerializeByValue(feature2.ID)
		assert.NilError(t, err)
		value, err := SerializeByValue("F")
		assert.NilError(t, err)
		key, err := SerializeRedisKey(revision.ID, "3215")
		assert.NilError(t, err)
		gender, err := store.HGet(ctx, key, field).Result()
		assert.NilError(t, err)
		assert.Equal(t, gender, value)
	}
}
