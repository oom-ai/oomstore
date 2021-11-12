package oomstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	mock_metadatav2 "github.com/oom-ai/oomstore/internal/database/metadatav2/mock_metadata"
	"github.com/oom-ai/oomstore/internal/database/online/mock_online"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
	"github.com/stretchr/testify/assert"
)

func TestGetOnlineFeatureValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	onlineStore := mock_online.NewMockStore(ctrl)
	metadatav2Store := mock_metadatav2.NewMockStore(ctrl)
	store := oomstore.NewOomStore(onlineStore, nil, metadatav2Store)

	entityName := "device"
	consistentFeatures := prepareFeatures(true, true)
	inconsistentFeatures := prepareFeatures(false, true)
	unavailableFeatures := prepareFeatures(true, false)

	testCases := []struct {
		description   string
		opt           types.GetOnlineFeatureValuesOpt
		entityName    *string
		features      typesv2.FeatureList
		expectedError error
		expected      types.FeatureValueMap
	}{
		{
			description: "no available features, return nil",
			opt: types.GetOnlineFeatureValuesOpt{
				FeatureNames: unavailableFeatures.Ids(),
				EntityKey:    "1234",
			},
			features:      unavailableFeatures,
			expectedError: nil,
			expected:      types.FeatureValueMap{},
		},
		{
			description: "inconsistent entity type, fail",
			opt: types.GetOnlineFeatureValuesOpt{
				FeatureNames: inconsistentFeatures.Ids(),
				EntityKey:    "1234",
			},
			features:      inconsistentFeatures,
			expectedError: fmt.Errorf("inconsistent entity type: %v", map[string]string{"device": "price", "user": "age"}),
			expected:      nil,
		},
		{
			description: "consistent entity type, succeed",
			opt: types.GetOnlineFeatureValuesOpt{
				FeatureNames: consistentFeatures.Ids(),
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
			metadatav2Store.EXPECT().ListFeature(gomock.Any(), metadatav2.ListFeatureOpt{FeatureIDs: tc.opt.FeatureNames}).Return(tc.features, nil)
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
	metadatav2Store := mock_metadatav2.NewMockStore(ctrl)

	store := oomstore.NewOomStore(onlineStore, nil, metadatav2Store)

	entityName := "device"
	consistentFeatures := prepareFeatures(true, true)
	inconsistentFeatures := prepareFeatures(false, true)
	unavailableFeatures := prepareFeatures(true, false)

	testCases := []struct {
		description   string
		opt           types.MultiGetOnlineFeatureValuesOpt
		entityName    *string
		features      typesv2.FeatureList
		expectedError error
		expected      types.FeatureDataSet
	}{
		{
			description: "no available features, return nil",
			opt: types.MultiGetOnlineFeatureValuesOpt{
				FeatureIDs: unavailableFeatures.Ids(),
				EntityKeys: []string{"1234", "1235"},
			},
			features:      unavailableFeatures,
			expectedError: nil,
			expected:      nil,
		},
		{
			description: "inconsistent entity type, fail",
			opt: types.MultiGetOnlineFeatureValuesOpt{
				FeatureIDs: inconsistentFeatures.Ids(),
				EntityKeys: []string{"1234", "1235"},
			},
			features:      inconsistentFeatures,
			expectedError: fmt.Errorf("inconsistent entity type: %v", map[string]string{"device": "price", "user": "age"}),
			expected:      nil,
		},
		{
			description: "consistent entity type, succeed",
			opt: types.MultiGetOnlineFeatureValuesOpt{
				FeatureIDs: consistentFeatures.Ids(),
				EntityKeys: []string{"1234", "1235"},
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
			metadatav2Store.EXPECT().ListFeature(gomock.Any(), metadatav2.ListFeatureOpt{FeatureIDs: tc.opt.FeatureIDs}).Return(tc.features, nil)
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

func prepareFeatures(isConsistent bool, isAvailable bool) typesv2.FeatureList {
	revision1 := int32(1)
	revision2 := int32(2)
	features := typesv2.FeatureList{
		{
			Name:        "model",
			DBValueType: "VARCHAR(32)",
			GroupID:     1,
			Group: &typesv2.FeatureGroup{
				EntityID:         1,
				OnlineRevisionID: &revision1,
				Category:         types.BatchFeatureCategory,
			},
		},
		{
			Name:        "price",
			DBValueType: "INT",
			GroupID:     2,
			Group: &typesv2.FeatureGroup{
				EntityID:         1,
				OnlineRevisionID: &revision2,
				Category:         types.BatchFeatureCategory,
			},
		},
		{
			Name:        "age",
			DBValueType: "INT",
			GroupID:     3,
			Group: &typesv2.FeatureGroup{
				EntityID:         2,
				OnlineRevisionID: &revision2,
				Category:         types.BatchFeatureCategory,
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
