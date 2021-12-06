package oomstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/metadata"
	"github.com/ethhte88/oomstore/internal/database/metadata/mock_metadata"
	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/internal/database/offline/mock_offline"
	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/internal/database/online/mock_online"
	"github.com/ethhte88/oomstore/pkg/oomstore"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	onlineStore := mock_online.NewMockStore(ctrl)
	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.TEST__New(onlineStore, offlineStore, metadataStore)

	features := types.FeatureList{
		{
			Name: "feature1",
		},
		{
			Name: "feature2",
		},
	}

	testCases := []struct {
		description   string
		opt           types.SyncOpt
		mockFn        func()
		expectedError error
	}{
		{
			description: "the specific revision was synced to the online store, won't do it again this time",
			opt: types.SyncOpt{
				RevisionID: 1,
			},
			expectedError: fmt.Errorf("the specific revision was synced to the online store, won't do it again this time"),
			mockFn: func() {
				metadataStore.EXPECT().GetRevision(ctx, 1).Return(&types.Revision{
					GroupID: 1,
					Group: &types.Group{
						ID:               1,
						OnlineRevisionID: intPtr(1),
					},
				}, nil)
			},
		},
		{
			description: "no previous revision, succeed",
			opt: types.SyncOpt{
				RevisionID: 1,
			},
			expectedError: nil,
			mockFn: func() {
				revision := buildRevision()
				metadataStore.EXPECT().GetRevision(ctx, 1).Return(revision, nil)
				metadataStore.EXPECT().GetGroupByName(ctx, "device_info").Return(&types.Group{
					Name: "device_info",
					ID:   1,
				}, nil)
				metadataStore.EXPECT().ListFeature(ctx, metadata.ListFeatureOpt{GroupID: &revision.Group.ID}).Return(features, nil)

				stream := make(chan types.ExportRecord)
				offlineStore.EXPECT().Export(ctx, offline.ExportOpt{
					DataTable:    "data-table-name",
					EntityName:   "device",
					FeatureNames: features.Names(),
				}).Return(stream, nil)

				onlineStore.EXPECT().Import(ctx, online.ImportOpt{
					FeatureList: features,
					Revision:    revision,
					Entity: &types.Entity{
						Name: "device",
					},
					ExportStream: stream,
				}).Return(nil)

				metadataStore.EXPECT().WithTransaction(ctx, gomock.Any()).Return(nil)
			},
		},
		{
			description: "purge previous revision, succeed",
			opt: types.SyncOpt{
				RevisionID: 1,
			},
			expectedError: nil,
			mockFn: func() {
				revision := buildRevision()
				revision.Group.OnlineRevisionID = intPtr(0)
				metadataStore.EXPECT().GetRevision(ctx, 1).Return(revision, nil)
				metadataStore.EXPECT().GetGroupByName(ctx, "device_info").Return(&types.Group{
					Name: "device_info",
					ID:   1,
				}, nil)
				metadataStore.EXPECT().ListFeature(ctx, metadata.ListFeatureOpt{GroupID: &revision.Group.ID}).Return(features, nil)

				stream := make(chan types.ExportRecord)
				offlineStore.EXPECT().Export(ctx, offline.ExportOpt{
					DataTable:    "data-table-name",
					EntityName:   "device",
					FeatureNames: features.Names(),
				}).Return(stream, nil)

				onlineStore.EXPECT().Import(ctx, online.ImportOpt{
					FeatureList: features,
					Revision:    revision,
					Entity: &types.Entity{
						Name: "device",
					},
					ExportStream: stream,
				}).Return(nil)

				metadataStore.EXPECT().WithTransaction(ctx, gomock.Any()).Return(nil)
				onlineStore.EXPECT().Purge(ctx, *revision.Group.OnlineRevisionID).Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().Refresh().Return(nil).AnyTimes()
			tc.mockFn()
			require.Equal(t, tc.expectedError, store.Sync(ctx, tc.opt))
		})
	}
}

func intPtr(i int) *int {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}

func buildRevision() *types.Revision {
	return &types.Revision{
		ID:      1,
		GroupID: 1,
		Group: &types.Group{
			Name:     "device_info",
			ID:       1,
			EntityID: 2,
			Entity: &types.Entity{
				Name: "device",
			},
			OnlineRevisionID: nil,
		},
		DataTable: "data-table-name",
	}
}
