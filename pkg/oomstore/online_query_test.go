package oomstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/mock_metadata"
	"github.com/oom-ai/oomstore/internal/database/online/mock_online"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

func TestOnlineGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	onlineStore := mock_online.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.TEST__New(onlineStore, nil, metadataStore)

	entityName := "device"
	consistentFeatures := prepareFeatures(true, true)
	inconsistentFeatures := prepareFeatures(false, true)
	unavailableFeatures := prepareFeatures(true, false)

	testCases := []struct {
		description   string
		opt           types.OnlineGetOpt
		entityName    *string
		features      types.FeatureList
		expectedError error
		expected      *types.FeatureValues
	}{
		{
			description: "no available features",
			opt: types.OnlineGetOpt{
				FeatureNames: unavailableFeatures.Names(),
				EntityKey:    "1234",
			},
			features:      unavailableFeatures,
			expectedError: nil,
			expected:      nil,
		},
		{
			description: "inconsistent entity type, fail",
			opt: types.OnlineGetOpt{
				FeatureNames: inconsistentFeatures.Names(),
				EntityKey:    "1234",
			},
			features:      inconsistentFeatures,
			expectedError: fmt.Errorf("expected 1 entity, got 2 entities"),
			expected:      nil,
		},
		{
			description: "consistent entity type, succeed",
			opt: types.OnlineGetOpt{
				FeatureNames: consistentFeatures.Names(),
				EntityKey:    "1234",
			},
			features:      consistentFeatures,
			entityName:    &entityName,
			expectedError: nil,
			expected: &types.FeatureValues{
				EntityName:   entityName,
				EntityKey:    "1234",
				FeatureNames: consistentFeatures.Names(),
				FeatureValueMap: map[string]interface{}{
					"price": int64(100),
					"model": "xiaomi",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().CacheListFeature(gomock.Any(), metadata.ListFeatureOpt{FeatureNames: &tc.opt.FeatureNames}).Return(tc.features)
			if tc.entityName != nil {
				onlineStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(dbutil.RowMap{
					"price": int64(100),
				}, nil)
				onlineStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(dbutil.RowMap{
					"model": "xiaomi",
				}, nil)
			}
			actual, err := store.OnlineGet(context.Background(), tc.opt)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestOnlineMultiGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	onlineStore := mock_online.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)

	store := oomstore.TEST__New(onlineStore, nil, metadataStore)

	entityName := "device"
	consistentFeatures := prepareFeatures(true, true)
	inconsistentFeatures := prepareFeatures(false, true)
	unavailableFeatures := prepareFeatures(true, false)

	testCases := []struct {
		description   string
		opt           types.OnlineMultiGetOpt
		entityName    *string
		features      types.FeatureList
		expectedError error
		expected      map[string]*types.FeatureValues
	}{
		{
			description: "no available features, return nil",
			opt: types.OnlineMultiGetOpt{
				FeatureNames: unavailableFeatures.Names(),
				EntityKeys:   []string{"1234", "1235"},
			},
			features:      unavailableFeatures,
			expectedError: nil,
			expected:      nil,
		},
		{
			description: "inconsistent entity type, fail",
			opt: types.OnlineMultiGetOpt{
				FeatureNames: inconsistentFeatures.Names(),
				EntityKeys:   []string{"1234", "1235"},
			},
			features:      inconsistentFeatures,
			expectedError: fmt.Errorf("expected 1 entity, got 2 entities"),
			expected:      nil,
		},
		{
			description: "consistent entity type, succeed",
			opt: types.OnlineMultiGetOpt{
				FeatureNames: consistentFeatures.Names(),
				EntityKeys:   []string{"1234", "1235"},
			},
			features:      consistentFeatures,
			entityName:    &entityName,
			expectedError: nil,
			expected: map[string]*types.FeatureValues{
				"1234": {
					EntityName:   "device",
					EntityKey:    "1234",
					FeatureNames: consistentFeatures.Names(),
					FeatureValueMap: map[string]interface{}{
						"model": "xiaomi",
						"price": int64(100),
					},
				},
				"1235": {
					EntityName:   "device",
					EntityKey:    "1235",
					FeatureNames: consistentFeatures.Names(),
					FeatureValueMap: map[string]interface{}{
						"model": "apple",
						"price": int64(200),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().CacheListFeature(gomock.Any(), metadata.ListFeatureOpt{FeatureNames: &tc.opt.FeatureNames}).Return(tc.features)
			if tc.entityName != nil {
				onlineStore.EXPECT().MultiGet(gomock.Any(), gomock.Any()).Return(map[string]dbutil.RowMap{
					"1234": {
						"price": int64(100),
					},
					"1235": {
						"price": int64(200),
					},
				}, nil)
				onlineStore.EXPECT().MultiGet(gomock.Any(), gomock.Any()).Return(map[string]dbutil.RowMap{
					"1234": {
						"model": "xiaomi",
					},
					"1235": {
						"model": "apple",
					},
				}, nil)
			}
			actual, err := store.OnlineMultiGet(context.Background(), tc.opt)
			if tc.expectedError != nil {
				assert.Error(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func prepareFeatures(isConsistent bool, isAvailable bool) types.FeatureList {
	revision1 := 1
	revision2 := 2
	entityDevice := &types.Entity{
		ID:   1,
		Name: "device",
	}
	entityUser := &types.Entity{
		ID:   2,
		Name: "user",
	}
	features := types.FeatureList{
		{
			Name:      "model",
			ValueType: types.String,
			GroupID:   1,
			Group: &types.Group{
				EntityID:         1,
				OnlineRevisionID: &revision1,
				Category:         types.BatchFeatureCategory,
				Entity:           entityDevice,
			},
		},
		{
			Name:      "price",
			ValueType: types.Int64,
			GroupID:   2,
			Group: &types.Group{
				EntityID:         1,
				OnlineRevisionID: &revision2,
				Category:         types.BatchFeatureCategory,
				Entity:           entityDevice,
			},
		},
		{
			Name:      "age",
			ValueType: types.Int64,
			GroupID:   3,
			Group: &types.Group{
				EntityID:         2,
				OnlineRevisionID: &revision2,
				Category:         types.BatchFeatureCategory,
				Entity:           entityUser,
			},
		},
	}
	if !isAvailable {
		for i := range features {
			features[i].Group.OnlineRevisionID = nil
			features[i].Group.Category = types.StreamFeatureCategory
		}
	}

	if !isConsistent {
		return features
	}
	return features[0:2]
}
