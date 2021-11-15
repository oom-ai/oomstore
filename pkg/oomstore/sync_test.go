package oomstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/mock_metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/mock_offline"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/mock_online"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func TestSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	onlineStore := mock_online.NewMockStore(ctrl)
	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.NewOomStore(onlineStore, offlineStore, metadataStore)
	ctx := context.Background()

	testCases := []struct {
		description   string
		opt           types.SyncOpt
		mockFn        func()
		expectedError error
	}{
		{
			description: "the specific revision was synced to the online store, won't do it again this time",
			opt: types.SyncOpt{
				RevisionId: 1,
			},
			expectedError: fmt.Errorf("the specific revision was synced to the online store, won't do it again this time"),
			mockFn: func() {
				metadataStore.EXPECT().
					GetRevision(ctx, int32(1)).
					Return(&types.Revision{
						GroupID: 1,
						Group: &types.FeatureGroup{
							ID:               1,
							OnlineRevisionID: int32Ptr(1),
						},
					}, nil)
			},
		},
		{
			description: "no previous revision, succeed",
			opt: types.SyncOpt{
				RevisionId: 1,
			},
			expectedError: nil,
			mockFn: func() {
				revision := &types.Revision{
					GroupID: 2,
					Group: &types.FeatureGroup{
						ID:       2,
						EntityID: 2,
						Entity: &types.Entity{
							Name: "entity-name",
						},
						OnlineRevisionID: nil,
					},
					DataTable: "data-table-name",
				}
				metadataStore.EXPECT().
					GetRevision(ctx, int32(1)).
					Return(revision, nil)

				features := types.FeatureList{
					{
						Name: "feature1",
					},
					{
						Name: "feature2",
					},
					{
						Name: "feature3",
					},
				}

				metadataStore.EXPECT().
					ListFeature(ctx, metadata.ListFeatureOpt{GroupID: &revision.Group.ID}).
					Return(features)

				stream := make(chan *types.RawFeatureValueRecord)
				offlineStore.EXPECT().
					Export(ctx, offline.ExportOpt{
						DataTable:    "data-table-name",
						EntityName:   "entity-name",
						FeatureNames: features.Names(),
					}).
					Return(stream, nil)

				onlineStore.EXPECT().
					Import(ctx, online.ImportOpt{
						FeatureList: features,
						Revision:    revision,
						Entity: &types.Entity{
							Name: "entity-name",
						},
						Stream: stream,
					}).
					Return(nil)

				metadataStore.EXPECT().
					UpdateFeatureGroup(ctx, metadata.UpdateFeatureGroupOpt{
						GroupID:             revision.GroupID,
						NewOnlineRevisionID: int32Ptr(revision.ID),
					}).
					Return(nil)

				metadataStore.EXPECT().
					UpdateRevision(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			description: "user-defined revision, succeed",
			opt: types.SyncOpt{
				RevisionId: 1,
			},
			expectedError: nil,
			mockFn: func() {
				revision := &types.Revision{
					GroupID: 2,
					Group: &types.FeatureGroup{
						ID:       2,
						EntityID: 2,
						Entity: &types.Entity{
							Name: "entity-name",
						},
						OnlineRevisionID: int32Ptr(100),
					},
					DataTable: "data-table-name",
				}
				metadataStore.EXPECT().
					GetRevision(ctx, int32(1)).
					Return(revision, nil)

				features := types.FeatureList{
					{
						Name: "feature1",
					},
					{
						Name: "feature2",
					},
					{
						Name: "feature3",
					},
				}

				metadataStore.EXPECT().
					ListFeature(ctx, metadata.ListFeatureOpt{GroupID: &revision.Group.ID}).
					Return(features)

				stream := make(chan *types.RawFeatureValueRecord)
				offlineStore.EXPECT().
					Export(ctx, offline.ExportOpt{
						DataTable:    "data-table-name",
						EntityName:   "entity-name",
						FeatureNames: features.Names(),
					}).
					Return(stream, nil)

				onlineStore.EXPECT().
					Import(ctx, online.ImportOpt{
						FeatureList: features,
						Revision:    revision,
						Entity: &types.Entity{
							Name: "entity-name",
						},
						Stream: stream,
					}).
					Return(nil)

				metadataStore.EXPECT().
					UpdateFeatureGroup(ctx, metadata.UpdateFeatureGroupOpt{
						GroupID:             revision.GroupID,
						NewOnlineRevisionID: int32Ptr(revision.ID),
					}).
					Return(nil)

				onlineStore.EXPECT().
					Purge(ctx, *revision.Group.OnlineRevisionID).
					Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc.mockFn()
			require.Equal(t, tc.expectedError, store.Sync(ctx, tc.opt))
		})
	}
}

func int16Ptr(i int16) *int16 {
	return &i
}

func int32Ptr(i int32) *int32 {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}
