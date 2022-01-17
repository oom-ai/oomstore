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

func TestChannelJoin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.TEST__New(nil, offlineStore, metadataStore)

	inconsistentFeatures := prepareFeatures(false, true)
	consistentFeatures := prepareFeatures(true, true)

	entity := types.Entity{
		Name: "device",
	}
	revisions := types.RevisionList{
		{
			Revision:      1,
			SnapshotTable: "device_basic_1",
		},
		{
			Revision:      20,
			SnapshotTable: "device_basic_20",
		},
	}

	validResult := prepareResult()
	entityRows := make(chan types.EntityRow)

	testCases := []struct {
		description   string
		opt           types.ChannelJoinOpt
		features      types.FeatureList
		entity        *types.Entity
		featureMap    map[string]types.FeatureList
		joined        *types.JoinResult
		expectedError error
		expected      *types.JoinResult
	}{
		{
			description: "no valid features, return nil",
			opt: types.ChannelJoinOpt{
				JoinFeatureNames: []string{},
				EntityRows:       entityRows,
			},
			features:      nil,
			expectedError: nil,
			expected:      prepareEmptyResult(),
		},
		{
			description: "inconsistent features, return nil",
			opt: types.ChannelJoinOpt{
				JoinFeatureNames: inconsistentFeatures.FullNames(),
				EntityRows:       entityRows,
			},
			features:      inconsistentFeatures,
			expectedError: fmt.Errorf("inconsistent entity type: %v", map[string]string{"device": "price", "user": "age"}),
			expected:      nil,
		},
		{
			description: "nil joined, return nil",
			opt: types.ChannelJoinOpt{
				JoinFeatureNames: consistentFeatures.FullNames(),
				EntityRows:       entityRows,
			},
			entity:   &entity,
			features: consistentFeatures,
			featureMap: map[string]types.FeatureList{
				"device_basic":    {consistentFeatures[0]},
				"device_advanced": {consistentFeatures[1]},
			},
			joined:        prepareEmptyResult(),
			expectedError: nil,
			expected:      prepareEmptyResult(),
		},
		{
			description: "consistent entity type, succeed",
			opt: types.ChannelJoinOpt{
				JoinFeatureNames: consistentFeatures.FullNames(),
				EntityRows:       entityRows,
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
			metadataStore.EXPECT().Refresh().Return(nil).AnyTimes()
			metadataStore.EXPECT().ListFeature(gomock.Any(), metadata.ListFeatureOpt{}).Return(tc.features, nil)
			if tc.entity != nil {
				for _, featureList := range tc.featureMap {
					metadataStore.EXPECT().ListRevision(gomock.Any(), &featureList[0].GroupID).Return(revisions, nil).AnyTimes()
				}
				offlineStore.EXPECT().Join(gomock.Any(), gomock.Any()).Return(tc.joined, nil)
			}

			actual, err := store.ChannelJoin(context.Background(), tc.opt)
			if tc.expectedError != nil {
				assert.Error(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.Header, actual.Header)
				assert.ObjectsAreEqual(extractValues(tc.expected.Data), extractValues(actual.Data))
			}
		})
	}
}

func prepareResult() *types.JoinResult {
	header := []string{"entity_key", "unix_milli", "model", "price"}
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
	if stream == nil {
		return values
	}
	for item := range stream {
		values = append(values, item)
	}
	return values
}

func prepareEmptyResult() *types.JoinResult {
	data := make(chan []interface{})
	defer close(data)
	return &types.JoinResult{
		Data: data,
	}
}
