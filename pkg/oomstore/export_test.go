package oomstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/internal/database/metadata/mock_metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/mock_offline"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestChannelExport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.TEST__New(nil, offlineStore, metadataStore)

	features := types.FeatureList{
		{
			Name:      "model",
			ValueType: types.String,
		},
		{
			Name:      "price",
			ValueType: types.Int64,
		},
	}

	testCases := []struct {
		description  string
		opt          types.ChannelExportOpt
		features     types.FeatureList
		exportStream <-chan types.ExportRecord
		exportError  <-chan error
		expected     [][]interface{}
	}{
		{
			description: "provide multiple features, should return multiple feature values",
			opt: types.ChannelExportOpt{
				FeatureFullNames: []string{"model", "price"},
				UnixMilli:        100,
			},
			features:     features,
			exportStream: prepareTwoFeatureStream(),
			expected:     [][]interface{}{{"1234", "xiaomi", int64(100)}, {"1235", "apple", int64(200)}, {"1236", "huawei", int64(300)}, {"1237", "oneplus", int64(240)}},
		},
		{
			description: "provide no features, should return empty values",
			opt: types.ChannelExportOpt{
				FeatureFullNames: []string{},
				UnixMilli:        100,
			},
			features:     features,
			exportStream: prepareEmptyStream(),
			exportError:  prepareExportError(),
			expected:     [][]interface{}{},
		},
		{
			description: "provide one feature, should return one feature values",
			opt: types.ChannelExportOpt{
				FeatureFullNames: []string{"price"},
				UnixMilli:        100,
			},
			features:     features[1:],
			exportStream: prepareOneFeatureStream(),
			expected:     [][]interface{}{{"1234", int64(100)}, {"1235", int64(200)}, {"1236", int64(300)}, {"1237", int64(240)}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			t.Skip()
			metadataStore.EXPECT().Refresh().Return(nil).AnyTimes()

			if len(tc.opt.FeatureFullNames) == 0 {
				metadataStore.EXPECT().GetGroupByName(ctx, "device_info").Return(&types.Group{
					Name: "device_info",
					ID:   1,
				}, nil)
			}
			metadataStore.EXPECT().ListFeature(gomock.Any(), gomock.Any()).Return(features, nil)

			offlineStore.EXPECT().Export(gomock.Any(), offline.ExportOpt{
				SnapshotTables: map[int]string{1: "device_info_10"},
				Features:       map[int]types.FeatureList{1: features},
				EntityName:     "device",
				Limit:          tc.opt.Limit,
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
