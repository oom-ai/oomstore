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
	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.NewOomStore(nil, offlineStore, metadataStore)

	testCases := []struct {
		description string
		opt         types.CreateFeatureOpt
		valueType   string
		group       types.FeatureGroup
		expectError bool
	}{
		{
			description: "create batch feature, succeed",
			opt: types.CreateFeatureOpt{
				FeatureName: "model",
				GroupID:     1,
				DBValueType: "VARCHAR(32)",
			},
			valueType: types.STRING,
			group: types.FeatureGroup{
				Name:     "device_info",
				Category: types.BatchFeatureCategory,
			},
			expectError: false,
		},
		{
			description: "create stream feature, fail",
			opt: types.CreateFeatureOpt{
				FeatureName: "model",
				GroupID:     1,
				DBValueType: "BIGINT",
			},
			valueType: types.INT64,
			group: types.FeatureGroup{
				Name:     "device_info",
				Category: types.StreamFeatureCategory,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().
				GetFeatureGroup(gomock.Any(), tc.opt.GroupID).
				Return(&tc.group, nil)

			if tc.group.Category == types.BatchFeatureCategory {
				metadataOpt := metadata.CreateFeatureOpt{
					CreateFeatureOpt: tc.opt,
					ValueType:        tc.valueType,
				}
				offlineStore.EXPECT().TypeTag(tc.opt.DBValueType).Return(tc.valueType, nil)
				metadataStore.EXPECT().CreateFeature(gomock.Any(), metadataOpt).Return(0, nil)
			}
			_, err := store.CreateBatchFeature(context.Background(), tc.opt)
			if tc.expectError {
				assert.Error(t, err, fmt.Errorf("expected batch feature group, got %s feature group", tc.group.Category))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
