package oomstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/mock_metadata"
	"github.com/oom-ai/oomstore/internal/database/offline/mock_offline"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestCreateBatchFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.TEST__New(nil, offlineStore, metadataStore)

	testCases := []struct {
		description string
		opt         types.CreateFeatureOpt
		valueType   string
		group       types.Group
		expectError bool
	}{
		{
			description: "create batch feature, succeed",
			opt: types.CreateFeatureOpt{
				FeatureName: "model",
				GroupName:   "device_info",
				DBValueType: "VARCHAR(32)",
			},
			valueType: types.STRING,
			group: types.Group{
				ID:       1,
				Name:     "device_info",
				Category: types.BatchFeatureCategory,
			},
			expectError: false,
		},
		{
			description: "create stream feature, fail",
			opt: types.CreateFeatureOpt{
				FeatureName: "model",
				GroupName:   "device_info",
				DBValueType: "BIGINT",
			},
			valueType: types.INT64,
			group: types.Group{
				ID:       1,
				Name:     "device_info",
				Category: types.StreamFeatureCategory,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().Refresh().Return(nil).AnyTimes()
			metadataStore.EXPECT().GetGroupByName(ctx, tc.opt.GroupName).Return(&tc.group, nil)

			if tc.group.Category == types.BatchFeatureCategory {
				metadataOpt := metadata.CreateFeatureOpt{
					FeatureName: tc.opt.FeatureName,
					GroupID:     tc.group.ID,
					DBValueType: tc.opt.DBValueType,
					ValueType:   tc.valueType,
					Description: tc.opt.Description,
				}
				offlineStore.EXPECT().TypeTag(tc.opt.DBValueType).Return(tc.valueType, nil)
				metadataStore.EXPECT().CreateFeature(ctx, metadataOpt).Return(0, nil)
			}
			_, err := store.CreateBatchFeature(ctx, tc.opt)
			if tc.expectError {
				assert.Error(t, err, fmt.Errorf("expected batch feature group, got %s feature group", tc.group.Category))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
