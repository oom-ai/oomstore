package oomstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata/mock_metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/mock_offline"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

func TestGetHistoricalFeatureValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)

	store := oomstore.NewOomStore(nil, offlineStore, metadataStore)
	streamFeatures := prepareFeatures(true, false)
	inconsistentFeatures := prepareFeatures(false, true)
	consistentFeatures := prepareFeatures(true, true)

	entity := types.Entity{
		Name:   "device",
		Length: 10,
	}
	revisionRanges := []*types.RevisionRange{
		{
			MinRevision: 1,
			MaxRevision: 20,
			DataTable:   "device_basic_1",
		},
	}
	entityRows, joinResult, result := prepareEntityRowsAndResult()

	testCases := []struct {
		description   string
		opt           types.GetHistoricalFeatureValuesOpt
		features      types.FeatureList
		entity        *types.Entity
		featureMap    map[string]types.FeatureList
		joinResultMap map[string]map[string]dbutil.RowMap
		expectedError error
		expected      []*types.EntityRowWithFeatures
	}{
		{
			description: "no valid features, return nil",
			opt: types.GetHistoricalFeatureValuesOpt{
				FeatureNames: streamFeatures.Names(),
				EntityRows:   entityRows,
			},
			features:      streamFeatures,
			expectedError: nil,
			expected:      nil,
		},
		{
			description: "inconsistent features, return nil",
			opt: types.GetHistoricalFeatureValuesOpt{
				FeatureNames: inconsistentFeatures.Names(),
				EntityRows:   entityRows,
			},
			features:      inconsistentFeatures,
			expectedError: fmt.Errorf("inconsistent entity type: %v", map[string]string{"device": "price", "user": "age"}),
			expected:      nil,
		},
		{
			description: "consistent entity type, succeed",
			opt: types.GetHistoricalFeatureValuesOpt{
				FeatureNames: consistentFeatures.Names(),
				EntityRows:   entityRows,
			},
			entity:   &entity,
			features: consistentFeatures,
			featureMap: map[string]types.FeatureList{
				"device_basic":    {consistentFeatures[0]},
				"device_advanced": {consistentFeatures[1]},
			},
			joinResultMap: joinResult,
			expectedError: nil,
			expected:      result,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().ListFeature(gomock.Any(), types.ListFeatureOpt{FeatureNames: tc.opt.FeatureNames}).Return(tc.features, nil)
			if tc.entity != nil {
				metadataStore.EXPECT().GetEntity(gomock.Any(), tc.entity.Name).Return(&entity, nil)
				for groupName, featureSlice := range tc.featureMap {
					metadataStore.EXPECT().BuildRevisionRanges(gomock.Any(), groupName).Return(revisionRanges, nil).AnyTimes()
					offlineStore.EXPECT().Join(gomock.Any(), offline.JoinOpt{
						Entity:         &entity,
						EntityRows:     tc.opt.EntityRows,
						RevisionRanges: revisionRanges,
						Features:       featureSlice,
					}).Return(tc.joinResultMap[groupName], nil)
				}
			}

			actual, err := store.GetHistoricalFeatureValues(context.Background(), tc.opt)
			if tc.expectedError != nil {
				assert.Error(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func prepareEntityRowsAndResult() ([]types.EntityRow, map[string]map[string]dbutil.RowMap, []*types.EntityRowWithFeatures) {
	entityRows := []types.EntityRow{
		{
			EntityKey: "1234",
			UnixTime:  10,
		},
		{
			EntityKey: "1234",
			UnixTime:  20,
		},
		{
			EntityKey: "1235",
			UnixTime:  15,
		},
		{
			EntityKey: "1236",
			UnixTime:  20,
		},
	}
	joinResult := map[string]map[string]dbutil.RowMap{
		"device_basic": {
			"1234,10": {
				"model":      "apple",
				"entity_key": "1234",
				"unix_time":  10,
			},
			"1234,20": {
				"model":      "oneplus",
				"entity_key": "1234",
				"unix_time":  20,
			},
			"1235,15": {
				"model":      "huawei",
				"entity_key": "1235",
				"unix_time":  15,
			},
		},
		"device_advanced": {
			"1234,10": {
				"price":      100,
				"entity_key": "1234",
				"unix_time":  10,
			},
			"1234,20": {
				"price":      120,
				"entity_key": "1234",
				"unix_time":  20,
			},
			"1235,15": {
				"price":      90,
				"entity_key": "1235",
				"unix_time":  15,
			},
		},
	}
	result := []*types.EntityRowWithFeatures{
		{
			EntityRow: types.EntityRow{
				EntityKey: "1234",
				UnixTime:  10,
			},
			FeatureValues: []types.FeatureKV{
				{
					FeatureName:  "model",
					FeatureValue: "apple",
				},
				{
					FeatureName:  "price",
					FeatureValue: 100,
				},
			},
		},
		{
			EntityRow: types.EntityRow{
				EntityKey: "1234",
				UnixTime:  20,
			},
			FeatureValues: []types.FeatureKV{
				{
					FeatureName:  "model",
					FeatureValue: "oneplus",
				},
				{
					FeatureName:  "price",
					FeatureValue: 120,
				},
			},
		},
		{
			EntityRow: types.EntityRow{
				EntityKey: "1235",
				UnixTime:  15,
			},
			FeatureValues: []types.FeatureKV{
				{
					FeatureName:  "model",
					FeatureValue: "huawei",
				},
				{
					FeatureName:  "price",
					FeatureValue: 90,
				},
			},
		},
		{
			EntityRow: types.EntityRow{
				EntityKey: "1236",
				UnixTime:  20,
			},
			FeatureValues: []types.FeatureKV{
				{
					FeatureName: "model",
				},
				{
					FeatureName: "price",
				},
			},
		},
	}
	return entityRows, joinResult, result
}
