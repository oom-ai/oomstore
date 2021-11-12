package oomstore_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	mock_metadatav2 "github.com/oom-ai/oomstore/internal/database/metadatav2/mock_metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/mock_offline"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
	"github.com/stretchr/testify/assert"
)

func TestExportFeatureValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	offlineStore := mock_offline.NewMockStore(ctrl)
	metadatav2Store := mock_metadatav2.NewMockStore(ctrl)

	store := oomstore.NewOomStore(nil, offlineStore, metadatav2Store)

	dataTable := "device_info_10"
	prevRevision := int64(5)
	prevDataTable := "device_info_5"
	features := typesv2.FeatureList{
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
				GroupID:      1,
				FeatureNames: []string{},
			},
			stream:   prepareTwoFeatureStream(),
			expected: [][]interface{}{{"1234", "xiaomi", int64(100)}, {"1235", "apple", int64(200)}, {"1236", "huawei", int64(300)}, {"1237", "oneplus", int64(240)}},
		},
		{
			description: "provide one feature",
			opt: types.ExportFeatureValuesOpt{
				GroupID:      1,
				FeatureNames: []string{"price"},
			},
			stream:   prepareOneFeatureStream(),
			expected: [][]interface{}{{"1234", int64(100)}, {"1235", int64(200)}, {"1236", int64(300)}, {"1237", int64(240)}},
		},
		{
			description: "provide revision",
			opt: types.ExportFeatureValuesOpt{
				GroupID:       1,
				FeatureNames:  []string{"price"},
				GroupRevision: &prevRevision,
			},
			stream:   prepareTwoFeatureStream(),
			expected: [][]interface{}{{"1234", "xiaomi", int64(100)}, {"1235", "apple", int64(200)}, {"1236", "huawei", int64(300)}, {"1237", "oneplus", int64(240)}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// mock database methods
			metadatav2Store.EXPECT().GetFeatureGroup(gomock.Any(), tc.opt.GroupID).Return(&typesv2.FeatureGroup{
				ID:       1,
				Name:     "device_info",
				EntityID: 1,
			}, nil)
			metadatav2Store.EXPECT().ListFeature(gomock.Any(), metadatav2.ListFeatureOpt{GroupID: &tc.opt.GroupID}).Return(features)

			dt := dataTable
			if tc.opt.GroupRevision != nil {
				dt = prevDataTable
				metadatav2Store.EXPECT().GetRevision(gomock.Any(), metadatav2.GetRevisionOpt{
					GroupID:  &tc.opt.GroupID,
					Revision: tc.opt.GroupRevision,
				}).Return(&typesv2.Revision{
					DataTable: prevDataTable,
				}, nil)
			}
			featureNames := tc.opt.FeatureNames
			if len(featureNames) == 0 {
				featureNames = features.Names()
			}

			offlineStore.EXPECT().Export(gomock.Any(), offline.ExportOpt{
				DataTable:    dt,
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
