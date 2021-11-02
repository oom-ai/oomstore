package oomstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata/mock_metadata"
	"github.com/oom-ai/oomstore/internal/database/online/mock_online"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

func TestGetOnlineFeatureValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	onlineStore := mock_online.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)

	store := oomstore.NewOomStore(onlineStore, nil, metadataStore)

	entityName := "device"
	consistentFeatures := prepareFeatures(true, true)
	inconsistentFeatures := prepareFeatures(false, true)
	unavailableFeatures := prepareFeatures(true, false)

	testCases := []struct {
		description   string
		opt           types.GetOnlineFeatureValuesOpt
		entityName    *string
		features      types.FeatureList
		expectedError error
		expected      types.FeatureValueMap
	}{
		{
			description: "no available features, return nil",
			opt: types.GetOnlineFeatureValuesOpt{
				FeatureNames: unavailableFeatures.Names(),
				EntityKey:    "1234",
			},
			features:      unavailableFeatures,
			expectedError: nil,
			expected:      types.FeatureValueMap{},
		},
		{
			description: "inconsistent entity type, fail",
			opt: types.GetOnlineFeatureValuesOpt{
				FeatureNames: inconsistentFeatures.Names(),
				EntityKey:    "1234",
			},
			features:      inconsistentFeatures,
			expectedError: fmt.Errorf("inconsistent entity type: %v", map[string]string{"device": "price", "user": "age"}),
			expected:      nil,
		},
		{
			description: "consistent entity type, succeed",
			opt: types.GetOnlineFeatureValuesOpt{
				FeatureNames: consistentFeatures.Names(),
				EntityKey:    "1234",
			},
			features:      consistentFeatures,
			entityName:    &entityName,
			expectedError: nil,
			expected: types.FeatureValueMap{
				"price": int64(100),
				"model": "xiaomi",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().ListFeature(gomock.Any(), types.ListFeatureOpt{FeatureNames: tc.opt.FeatureNames}).Return(tc.features, nil)
			if tc.entityName != nil {
				onlineStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(dbutil.RowMap{
					"price": int64(100),
				}, nil)
				onlineStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(dbutil.RowMap{
					"model": "xiaomi",
				}, nil)
			}
			actual, err := store.GetOnlineFeatureValues(context.Background(), tc.opt)
			if tc.expectedError != nil {
				assert.Error(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMultiGetOnlineFeatureValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	onlineStore := mock_online.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)

	store := oomstore.NewOomStore(onlineStore, nil, metadataStore)

	entityName := "device"
	consistentFeatures := prepareFeatures(true, true)
	inconsistentFeatures := prepareFeatures(false, true)
	unavailableFeatures := prepareFeatures(true, false)

	testCases := []struct {
		description   string
		opt           types.MultiGetOnlineFeatureValuesOpt
		entityName    *string
		features      types.FeatureList
		expectedError error
		expected      types.FeatureDataSet
	}{
		{
			description: "no available features, return nil",
			opt: types.MultiGetOnlineFeatureValuesOpt{
				FeatureNames: unavailableFeatures.Names(),
				EntityKeys:   []string{"1234", "1235"},
			},
			features:      unavailableFeatures,
			expectedError: nil,
			expected:      nil,
		},
		{
			description: "inconsistent entity type, fail",
			opt: types.MultiGetOnlineFeatureValuesOpt{
				FeatureNames: inconsistentFeatures.Names(),
				EntityKeys:   []string{"1234", "1235"},
			},
			features:      inconsistentFeatures,
			expectedError: fmt.Errorf("inconsistent entity type: %v", map[string]string{"device": "price", "user": "age"}),
			expected:      nil,
		},
		{
			description: "consistent entity type, succeed",
			opt: types.MultiGetOnlineFeatureValuesOpt{
				FeatureNames: consistentFeatures.Names(),
				EntityKeys:   []string{"1234", "1235"},
			},
			features:      consistentFeatures,
			entityName:    &entityName,
			expectedError: nil,
			expected: types.FeatureDataSet{
				"1234": []types.FeatureKV{
					{
						FeatureName:  "model",
						FeatureValue: "xiaomi",
					},
					{
						FeatureName:  "price",
						FeatureValue: int64(100),
					},
				},
				"1235": []types.FeatureKV{
					{
						FeatureName:  "model",
						FeatureValue: "apple",
					},
					{
						FeatureName:  "price",
						FeatureValue: int64(200),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().ListFeature(gomock.Any(), types.ListFeatureOpt{FeatureNames: tc.opt.FeatureNames}).Return(tc.features, nil)
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
			actual, err := store.MultiGetOnlineFeatureValues(context.Background(), tc.opt)
			if tc.expectedError != nil {
				assert.Error(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func prepareFeatures(isConsistent bool, available bool) types.FeatureList {
	revision1 := int32(1)
	revision2 := int32(2)
	features := types.FeatureList{
		{
			Name:             "model",
			DBValueType:      "VARCHAR(32)",
			EntityName:       "device",
			OnlineRevisionID: &revision1,
		},
		{
			Name:             "price",
			DBValueType:      "INT",
			EntityName:       "device",
			OnlineRevisionID: &revision2,
		},
		{
			Name:             "age",
			DBValueType:      "INT",
			EntityName:       "user",
			OnlineRevisionID: &revision2,
		},
	}
	if !available {
		for i := range features {
			features[i].OnlineRevisionID = nil
		}
	}

	if !isConsistent {
		return features
	}
	return features[0:2]
}
