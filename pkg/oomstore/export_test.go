package oomstore_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/internal/database/metadata/mock_metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/mock_offline"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestExportFeatureValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.NewOomStore(nil, offlineStore, metadataStore)

	revisionID := 5
	features := types.FeatureList{
		{
			Name:        "model",
			DBValueType: "VARCHAR(32)",
		},
		{
			Name:        "price",
			DBValueType: "INT",
		},
	}

	testCases := []struct {
		description string
		opt         types.ExportFeatureValuesOpt
		stream      <-chan *types.RawFeatureValueRecord
		expected    [][]interface{}
	}{
		{
			description: "no features",
			opt: types.ExportFeatureValuesOpt{
				FeatureNames: []string{},
				RevisionID:   revisionID,
			},
			stream:   prepareTwoFeatureStream(),
			expected: [][]interface{}{{"1234", "xiaomi", int64(100)}, {"1235", "apple", int64(200)}, {"1236", "huawei", int64(300)}, {"1237", "oneplus", int64(240)}},
		},
		{
			description: "provide one feature",
			opt: types.ExportFeatureValuesOpt{
				FeatureNames: []string{"price"},
				RevisionID:   revisionID,
			},
			stream:   prepareOneFeatureStream(),
			expected: [][]interface{}{{"1234", int64(100)}, {"1235", int64(200)}, {"1236", int64(300)}, {"1237", int64(240)}},
		},
		{
			description: "provide revision",
			opt: types.ExportFeatureValuesOpt{
				FeatureNames: []string{"price"},
				RevisionID:   revisionID,
			},
			stream:   prepareTwoFeatureStream(),
			expected: [][]interface{}{{"1234", "xiaomi", int64(100)}, {"1235", "apple", int64(200)}, {"1236", "huawei", int64(300)}, {"1237", "oneplus", int64(240)}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			revision := types.Revision{
				ID:        1,
				GroupID:   1,
				DataTable: "device_info_10",
				Group: &types.FeatureGroup{
					Name:     "device_info",
					ID:       1,
					EntityID: 1,
					Entity:   &types.Entity{Name: "device"},
				},
			}

			metadataStore.EXPECT().GetRevision(ctx, tc.opt.RevisionID).Return(&revision, nil)
			metadataStore.EXPECT().GetFeatureGroupByName(ctx, "device_info").Return(&types.FeatureGroup{
				Name: "device_info",
				ID:   1,
			}, nil)
			metadataStore.EXPECT().ListFeature(gomock.Any(), gomock.Any()).Return(features)

			featureNames := tc.opt.FeatureNames
			if len(featureNames) == 0 {
				featureNames = features.Names()
			}

			offlineStore.EXPECT().Export(gomock.Any(), offline.ExportOpt{
				DataTable:    "device_info_10",
				EntityName:   "device",
				FeatureNames: featureNames,
				Limit:        tc.opt.Limit,
			}).Return(tc.stream, nil)

			// execute and compare results
			_, actual, err := store.ExportFeatureValues(context.Background(), tc.opt)
			assert.NoError(t, err)
			values := make([][]interface{}, 0)
			for ele := range actual {
				values = append(values, ele.Record)
				assert.NoError(t, ele.Error)
			}
			assert.Equal(t, tc.expected, values)
		})
	}
}

func prepareTwoFeatureStream() chan *types.RawFeatureValueRecord {
	stream := make(chan *types.RawFeatureValueRecord)
	go func() {
		defer close(stream)
		stream <- &types.RawFeatureValueRecord{
			Record: []interface{}{"1234", "xiaomi", int64(100)},
		}
		stream <- &types.RawFeatureValueRecord{
			Record: []interface{}{"1235", "apple", int64(200)},
		}
		stream <- &types.RawFeatureValueRecord{
			Record: []interface{}{"1236", "huawei", int64(300)},
		}
		stream <- &types.RawFeatureValueRecord{
			Record: []interface{}{"1237", "oneplus", int64(240)},
		}
	}()
	return stream
}

func prepareOneFeatureStream() chan *types.RawFeatureValueRecord {
	stream := make(chan *types.RawFeatureValueRecord)
	go func() {
		defer close(stream)
		stream <- &types.RawFeatureValueRecord{
			Record: []interface{}{"1234", int64(100)},
		}
		stream <- &types.RawFeatureValueRecord{
			Record: []interface{}{"1235", int64(200)},
		}
		stream <- &types.RawFeatureValueRecord{
			Record: []interface{}{"1236", int64(300)},
		}
		stream <- &types.RawFeatureValueRecord{
			Record: []interface{}{"1237", int64(240)},
		}
	}()
	return stream
}
