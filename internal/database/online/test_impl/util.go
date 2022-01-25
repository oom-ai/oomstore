package test_impl

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
)

type PrepareStoreFn func(*testing.T) (context.Context, online.Store)

type DestroyStoreFn func()

type Sample struct {
	Entity   types.Entity
	Group    types.Group
	Features types.FeatureList
	Revision types.Revision
	Data     []types.ExportRecord
}

var SampleSmall Sample
var SampleMedium Sample

func init() {
	rand.Seed(time.Now().UnixNano())
	entity := types.Entity{ID: 5, Name: "user"}
	group1 := types.Group{ID: 1, Category: types.CategoryBatch, Entity: &entity}
	group2 := types.Group{ID: 2, Category: types.CategoryBatch, Entity: &entity}
	{
		SampleSmall = Sample{
			Entity: entity,
			Group:  group1,
			Features: types.FeatureList{
				&types.Feature{
					ID:        1,
					Name:      "age",
					GroupID:   1,
					Group:     &group1,
					ValueType: types.Int64,
				},
				&types.Feature{
					ID:        2,
					Name:      "gender",
					GroupID:   1,
					Group:     &group1,
					ValueType: types.String,
				},
				&types.Feature{
					ID:        3,
					Name:      "account",
					GroupID:   1,
					Group:     &group1,
					ValueType: types.Float64,
				},
				&types.Feature{
					ID:        4,
					Name:      "is_active",
					GroupID:   1,
					Group:     &group1,
					ValueType: types.Bool,
				},
				&types.Feature{
					ID:        5,
					Name:      "register_time",
					GroupID:   1,
					Group:     &group1,
					ValueType: types.Time,
				},
			},
			Revision: types.Revision{
				ID:      3,
				GroupID: 1,
				Group:   &group1,
			},
			Data: []types.ExportRecord{
				[]interface{}{"3215", int64(18), "F", 1.1, true, time.Now()},
				[]interface{}{"3216", int64(29), nil, 2.0, false, time.Now()},
				[]interface{}{"3217", int64(44), "M", 3.1, true, time.Now()},
			},
		}

	}

	{
		features := types.FeatureList{
			&types.Feature{
				ID:        2,
				Name:      "charge",
				GroupID:   2,
				Group:     &group2,
				ValueType: types.Float64,
			},
		}

		revision := types.Revision{ID: 9, GroupID: 2, Group: &group2}
		var data []types.ExportRecord

		for i := 0; i < 100; i++ {
			record := []interface{}{dbutil.RandString(10), rand.Float64()}
			data = append(data, record)
		}
		SampleMedium = Sample{entity, group1, features, revision, data}
	}
}

func importSample(t *testing.T, ctx context.Context, store online.Store, samples ...*Sample) {
	for _, sample := range samples {
		stream := make(chan types.ExportRecord)
		go func(sample *Sample) {
			defer close(stream)
			for i := range sample.Data {
				stream <- sample.Data[i]
			}
		}(sample)

		err := store.Import(ctx, online.ImportOpt{
			Group:        sample.Group,
			Features:     sample.Features,
			RevisionID:   &sample.Revision.ID,
			ExportStream: stream,
		})
		require.NoError(t, err)
	}
}

func compareFeatureValue(t *testing.T, expected, actual interface{}, valueType types.ValueType) {
	if valueType == types.Time {
		expected, ok := expected.(time.Time)
		require.Equal(t, true, ok)

		actual, ok := actual.(time.Time)
		require.Equal(t, true, ok)

		if expected.Location() == actual.Location() {
			assert.Equal(t, expected.Local().Unix(), actual.Local().Unix())
		}
	} else {
		assert.Equal(t, expected, actual)
	}
}
