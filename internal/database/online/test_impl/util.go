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
var SampleStream Sample

func init() {
	rand.Seed(time.Now().UnixNano())
	entity := types.Entity{ID: 1, Name: "user"}
	group1 := types.Group{ID: 1, Name: "group1", Category: types.CategoryBatch, Entity: &entity}
	group2 := types.Group{ID: 2, Name: "group2", Category: types.CategoryBatch, Entity: &entity}
	group3 := types.Group{ID: 3, Name: "user_clicks", Category: types.CategoryStream, Entity: &entity}

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
			{
				Record: []interface{}{"3215", int64(18), "F", 1.1, true, time.Now()},
			},
			{
				Record: []interface{}{"3216", int64(29), nil, 2.0, false, time.Now()},
			},
			{
				Record: []interface{}{"3217", int64(44), "M", 3.1, true, time.Now()},
			},
		},
	}

	var data []types.ExportRecord
	for i := 0; i < 100; i++ {
		record := []interface{}{dbutil.RandString(10), rand.Float64()}
		data = append(data, types.ExportRecord{Record: record, Error: nil})
	}
	SampleMedium = Sample{
		Entity: entity,
		Group:  group2,
		Features: types.FeatureList{
			&types.Feature{
				ID:        6,
				Name:      "charge",
				GroupID:   2,
				Group:     &group2,
				ValueType: types.Float64,
			},
		},
		Revision: types.Revision{
			ID:      9,
			GroupID: 2,
			Group:   &group2,
		},
		Data: data,
	}

	SampleStream = Sample{
		Entity: entity,
		Group:  group3,
		Features: types.FeatureList{
			&types.Feature{
				ID:        7,
				Name:      "amount",
				ValueType: types.Int64,
				GroupID:   3,
				Group:     &group3,
			},
			&types.Feature{
				ID:        8,
				Name:      "last_10_click_post",
				ValueType: types.String,
				GroupID:   3,
				Group:     &group3,
			},
		},
		Data: []types.ExportRecord{
			{
				Record: []interface{}{"3215", int64(1), "1,2,3,4"},
			},
			{
				Record: []interface{}{"3216", int64(2), "2,3,4,5"},
			},
			{
				Record: []interface{}{"3217", int64(3), "3,4,5,6"},
			},
			{
				Record: []interface{}{"3218", int64(4), "4,5,6,7"},
			},
		},
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

		opt := online.ImportOpt{
			Group:        sample.Group,
			Features:     sample.Features,
			ExportStream: stream,
		}
		if sample.Group.Category == types.CategoryBatch {
			opt.RevisionID = &sample.Revision.ID
		}
		err := store.Import(ctx, opt)
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
