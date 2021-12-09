package test_impl

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"

	"github.com/ethhte88/oomstore/internal/database/online"
)

type PrepareStoreFn func() (context.Context, online.Store)

type Sample struct {
	Features types.FeatureList
	Revision *types.Revision
	Entity   *types.Entity
	Data     []types.ExportRecord
}

var SampleSmall Sample
var SampleMedium Sample

func init() {
	rand.Seed(time.Now().UnixNano())

	{
		SampleSmall = Sample{
			Features: types.FeatureList{
				&types.Feature{
					ID:          1,
					Name:        "age",
					GroupID:     1,
					ValueType:   types.INT64,
					DBValueType: "int",
				},
				&types.Feature{
					ID:          2,
					Name:        "gender",
					GroupID:     1,
					ValueType:   types.STRING,
					DBValueType: "varchar(1)",
				},
				&types.Feature{
					ID:          3,
					Name:        "account",
					GroupID:     1,
					ValueType:   types.FLOAT64,
					DBValueType: "double",
				},
				&types.Feature{
					ID:          3,
					Name:        "is_active",
					GroupID:     1,
					ValueType:   types.BOOL,
					DBValueType: "bool",
				},
			},
			Revision: &types.Revision{ID: 3, GroupID: 1},
			Entity:   &types.Entity{ID: 5, Name: "user", Length: 4},
			Data: []types.ExportRecord{
				[]interface{}{"3215", int64(18), "F", 1.1, true},
				[]interface{}{"3216", int64(29), nil, 2.0, false},
				[]interface{}{"3217", int64(44), "M", 3.1, true},
			},
		}

	}

	{
		features := types.FeatureList{
			&types.Feature{
				ID:          2,
				Name:        "charge",
				GroupID:     2,
				ValueType:   types.FLOAT64,
				DBValueType: "float",
			},
		}

		revision := &types.Revision{ID: 9, GroupID: 2}
		entity := &types.Entity{ID: 5, Name: "user", Length: 5}
		var data []types.ExportRecord

		for i := 0; i < 1000; i++ {
			record := []interface{}{RandString(entity.Length), rand.Float64()}
			data = append(data, record)
		}
		SampleMedium = Sample{features, revision, entity, data}
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
			FeatureList:  sample.Features,
			Revision:     sample.Revision,
			Entity:       sample.Entity,
			ExportStream: stream,
		})
		require.NoError(t, err)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
