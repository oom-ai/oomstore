package oomstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/ethhte88/oomstore/internal/database/metadata/mock_metadata"
	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/internal/database/offline/mock_offline"
	"github.com/ethhte88/oomstore/pkg/oomstore"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func TestChannelExport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.TEST__New(nil, offlineStore, metadataStore)

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
		description  string
		opt          types.ChannelExportOpt
		exportStream <-chan types.ExportRecord
		exportError  <-chan error
		expected     [][]interface{}
	}{
		{
			description: "no features",
			opt: types.ChannelExportOpt{
				FeatureNames: []string{},
				RevisionID:   revisionID,
			},
			exportStream: prepareTwoFeatureStream(),
			expected:     [][]interface{}{{"1234", "xiaomi", int64(100)}, {"1235", "apple", int64(200)}, {"1236", "huawei", int64(300)}, {"1237", "oneplus", int64(240)}},
		},
		{
			description: "provide one feature",
			opt: types.ChannelExportOpt{
				FeatureNames: []string{"price"},
				RevisionID:   revisionID,
			},
			exportStream: prepareOneFeatureStream(),
			expected:     [][]interface{}{{"1234", int64(100)}, {"1235", int64(200)}, {"1236", int64(300)}, {"1237", int64(240)}},
		},
		{
			description: "provide revision",
			opt: types.ChannelExportOpt{
				FeatureNames: []string{"price"},
				RevisionID:   revisionID,
			},
			exportStream: prepareTwoFeatureStream(),
			expected:     [][]interface{}{{"1234", "xiaomi", int64(100)}, {"1235", "apple", int64(200)}, {"1236", "huawei", int64(300)}, {"1237", "oneplus", int64(240)}},
		},
		{
			description: "empty stream",
			opt: types.ChannelExportOpt{
				FeatureNames: []string{"price"},
				RevisionID:   revisionID,
			},
			exportStream: prepareEmptyStream(),
			exportError:  prepareExportError(),
			expected:     [][]interface{}{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().Refresh().Return(nil).AnyTimes()
			revision := types.Revision{
				ID:        1,
				GroupID:   1,
				DataTable: "device_info_10",
				Group: &types.Group{
					Name:     "device_info",
					ID:       1,
					EntityID: 1,
					Entity:   &types.Entity{Name: "device"},
				},
			}

			metadataStore.EXPECT().GetRevision(ctx, tc.opt.RevisionID).Return(&revision, nil)
			metadataStore.EXPECT().GetGroupByName(ctx, "device_info").Return(&types.Group{
				Name: "device_info",
				ID:   1,
			}, nil)
			metadataStore.EXPECT().ListFeature(gomock.Any(), gomock.Any()).Return(features, nil)

			featureNames := tc.opt.FeatureNames
			if len(featureNames) == 0 {
				featureNames = features.Names()
			}

			offlineStore.EXPECT().Export(gomock.Any(), offline.ExportOpt{
				DataTable:    "device_info_10",
				EntityName:   "device",
				FeatureNames: featureNames,
				Limit:        tc.opt.Limit,
			}).Return(tc.exportStream, tc.exportError)

			// execute and compare results
			actual, err := store.ChannelExport(context.Background(), tc.opt)
			assert.NoError(t, err)
			values := make([][]interface{}, 0)
			for row := range actual.Data {
				values = append(values, row)
			}
			if tc.exportError != nil {
				assert.Error(t, actual.CheckStreamError())
			} else {
				assert.NoError(t, actual.CheckStreamError())
			}
			assert.Equal(t, tc.expected, values)
		})
	}
}

func prepareTwoFeatureStream() chan types.ExportRecord {
	stream := make(chan types.ExportRecord)
	go func() {
		defer close(stream)
		stream <- []interface{}{"1234", "xiaomi", int64(100)}
		stream <- []interface{}{"1235", "apple", int64(200)}
		stream <- []interface{}{"1236", "huawei", int64(300)}
		stream <- []interface{}{"1237", "oneplus", int64(240)}
	}()
	return stream
}

func prepareOneFeatureStream() chan types.ExportRecord {
	stream := make(chan types.ExportRecord)
	go func() {
		defer close(stream)
		stream <- []interface{}{"1234", int64(100)}
		stream <- []interface{}{"1235", int64(200)}
		stream <- []interface{}{"1236", int64(300)}
		stream <- []interface{}{"1237", int64(240)}
	}()
	return stream
}

func prepareEmptyStream() chan types.ExportRecord {
	stream := make(chan types.ExportRecord)
	go func() {
		defer close(stream)
	}()
	return stream
}

func prepareExportError() <-chan error {
	err := make(chan error, 1)
	go func() {
		defer close(err)
		err <- fmt.Errorf("error")
	}()
	return err
}
