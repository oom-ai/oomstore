package oomstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/mock_metadata"
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
	store := oomstore.TEST__New(nil, offlineStore, metadataStore)

	streamFeatures := prepareFeatures(true, false)
	inconsistentFeatures := prepareFeatures(false, true)
	consistentFeatures := prepareFeatures(true, true)

	entity := types.Entity{
		Name:   "device",
		Length: 10,
	}
	revisions := types.RevisionList{
		{
			Revision:  1,
			DataTable: "device_basic_1",
		},
		{
			Revision:  20,
			DataTable: "device_basic_20",
		},
	}

	var emptyResult *types.JoinResult
	validResult := prepareResult()
	entityRows := make(chan types.EntityRow)

	testCases := []struct {
		description   string
		opt           types.JoinOpt
		features      types.FeatureList
		entity        *types.Entity
		featureMap    map[string]types.FeatureList
		joined        *types.JoinResult
		expectedError error
		expected      *types.JoinResult
	}{
		{
			description: "no valid features, return nil",
			opt: types.JoinOpt{
				FeatureIDs: streamFeatures.IDs(),
				EntityRows: entityRows,
			},
			features:      streamFeatures,
			expectedError: nil,
			expected:      emptyResult,
		},
		{
			description: "inconsistent features, return nil",
			opt: types.JoinOpt{
				FeatureIDs: inconsistentFeatures.IDs(),
				EntityRows: entityRows,
			},
			features:      inconsistentFeatures,
			expectedError: fmt.Errorf("inconsistent entity type: %v", map[string]string{"device": "price", "user": "age"}),
			expected:      nil,
		},
		{
			description: "nil joined, return nil",
			opt: types.JoinOpt{
				FeatureIDs: consistentFeatures.IDs(),
				EntityRows: entityRows,
			},
			entity:   &entity,
			features: consistentFeatures,
			featureMap: map[string]types.FeatureList{
				"device_basic":    {consistentFeatures[0]},
				"device_advanced": {consistentFeatures[1]},
			},
			joined:        nil,
			expectedError: nil,
			expected:      emptyResult,
		},
		{
			description: "consistent entity type, succeed",
			opt: types.JoinOpt{
				FeatureIDs: consistentFeatures.IDs(),
				EntityRows: entityRows,
			},
			entity:   &entity,
			features: consistentFeatures,
			featureMap: map[string]types.FeatureList{
				"device_basic":    {consistentFeatures[0]},
				"device_advanced": {consistentFeatures[1]},
			},
			joined:        validResult,
			expectedError: nil,
			expected:      validResult,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().ListFeature(gomock.Any(), metadata.ListFeatureOpt{FeatureIDs: &tc.opt.FeatureIDs}).Return(tc.features)
			if tc.entity != nil {
				for _, featureList := range tc.featureMap {
					metadataStore.EXPECT().ListRevision(gomock.Any(), metadata.ListRevisionOpt{GroupID: &featureList[0].GroupID}).Return(revisions).AnyTimes()
				}
				offlineStore.EXPECT().Join(gomock.Any(), gomock.Any()).Return(tc.joined, nil)
			}

			actual, err := store.Join(context.Background(), tc.opt)
			if tc.expectedError != nil {
				assert.Error(t, err, tc.expectedError)
			} else if tc.expected == nil {
				assert.Equal(t, tc.expected, actual)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.Header, actual.Header)
				assert.ObjectsAreEqual(extractValues(tc.expected.Data), extractValues(actual.Data))
			}
		})
	}
}

func prepareResult() *types.JoinResult {
	header := []string{"entity_key", "unix_time", "model", "price"}
	data := make(chan []interface{})
	go func() {
		defer close(data)
		data <- []interface{}{"1234", 10, "apple", 100}
		data <- []interface{}{"1234", 20, "oneplus", 120}
		data <- []interface{}{"1235", 15, "galaxy", 90}
	}()

	return &types.JoinResult{
		Header: header,
		Data:   data,
	}
}

func extractValues(stream <-chan []interface{}) [][]interface{} {
	values := make([][]interface{}, 0)
	for item := range stream {
		values = append(values, item)
	}
	return values
}
