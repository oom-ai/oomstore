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

func TestCreateFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.TEST__New(nil, offlineStore, metadataStore)

	testCases := []struct {
		description string
		opt         types.CreateFeatureOpt
		valueType   types.ValueType
		group       types.Group
	}{
		{
			description: "create batch feature, succeed",
			opt: types.CreateFeatureOpt{
				FeatureName: "model",
				GroupName:   "device_info",
				ValueType:   types.String,
			},
			valueType: types.String,
			group: types.Group{
				ID:       1,
				Name:     "device_info",
				Category: types.CategoryBatch,
			},
		},
		{
			description: "create stream feature, succeed",
			opt: types.CreateFeatureOpt{
				FeatureName: "model",
				GroupName:   "device_info",
				ValueType:   types.Int64,
			},
			valueType: types.Int64,
			group: types.Group{
				ID:       1,
				Name:     "device_info",
				Category: types.CategoryStream,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().Refresh().Return(nil).AnyTimes()
			metadataStore.EXPECT().GetGroupByName(ctx, tc.opt.GroupName).Return(&tc.group, nil)

			metadataOpt := metadata.CreateFeatureOpt{
				FeatureName: tc.opt.FeatureName,
				FullName:    fmt.Sprintf("%s.%s", tc.opt.GroupName, tc.opt.FeatureName),
				GroupID:     tc.group.ID,
				ValueType:   tc.valueType,
				Description: tc.opt.Description,
			}
			metadataStore.EXPECT().CreateFeature(ctx, metadataOpt).Return(0, nil)

			_, err := store.CreateFeature(ctx, tc.opt)
			assert.NoError(t, err)
		})
	}
}
