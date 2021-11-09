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
	"github.com/stretchr/testify/assert"
)

func TestMaterialize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	onlineStore := mock_online.NewMockStore(ctrl)
	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)

	store := oomstore.NewOomStore(onlineStore, offlineStore, metadataStore)
	ctx := context.Background()

	revision1 := prepareRevision(1, 10)
	revision2 := prepareRevision(2, 20)
	revision3 := prepareRevision(3, 30)
	entity := types.Entity{
		Name: "device",
	}
	features := prepareFeatures(true, true)
	stream := make(chan *types.RawFeatureValueRecord)

	testCases := []struct {
		description      string
		opt              types.MaterializeOpt
		group            types.FeatureGroup
		revision         types.Revision
		previousRevision types.Revision
		expectedError    error
	}{
		{
			description: "user-defined revision, succeed",
			opt: types.MaterializeOpt{
				GroupName:     "device_info",
				GroupRevision: int64Ptr(10),
			},
			group:            prepareGroup(int32Ptr(2)),
			revision:         revision1,
			previousRevision: revision2,
			expectedError:    nil,
		},
		{
			description: "latest revision, succeed",
			opt: types.MaterializeOpt{
				GroupName: "device_info",
			},
			group:            prepareGroup(int32Ptr(2)),
			revision:         revision3,
			previousRevision: revision2,
			expectedError:    nil,
		},
		{
			description: "no previous revision, succeed",
			opt: types.MaterializeOpt{
				GroupName: "device_info",
			},
			group:         prepareGroup(nil),
			revision:      revision1,
			expectedError: nil,
		},
		{
			description: "already in latest revision, fail",
			opt: types.MaterializeOpt{
				GroupName: "device_info",
			},
			group:         prepareGroup(int32Ptr(3)),
			revision:      revision3,
			expectedError: fmt.Errorf("online store already in the latest revision"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().GetFeatureGroup(ctx, tc.opt.GroupName).Return(&tc.group, nil)
			metadataStore.EXPECT().GetEntity(ctx, entity.Name).Return(&entity, nil)
			metadataStore.EXPECT().ListFeature(ctx, types.ListFeatureOpt{GroupName: &tc.opt.GroupName}).Return(features, nil)
			if tc.opt.GroupRevision != nil {
				metadataStore.EXPECT().GetRevision(ctx, metadata.GetRevisionOpt{
					GroupName: &tc.opt.GroupName,
					Revision:  tc.opt.GroupRevision,
				}).Return(&tc.revision, nil)
			} else {
				metadataStore.EXPECT().GetLatestRevision(ctx, tc.opt.GroupName).Return(&tc.revision, nil)
			}
			if tc.expectedError == nil {
				offlineStore.EXPECT().Export(ctx, offline.ExportOpt{
					DataTable:    tc.revision.DataTable,
					EntityName:   tc.group.EntityName,
					FeatureNames: features.Names(),
				}).Return(stream, nil)
				onlineStore.EXPECT().Import(ctx, online.ImportOpt{
					Features: features,
					Revision: &tc.revision,
					Entity:   &entity,
					Stream:   stream,
				})
				if tc.group.OnlineRevisionID != nil {
					metadataStore.EXPECT().GetRevision(ctx, metadata.GetRevisionOpt{
						RevisionId: tc.group.OnlineRevisionID,
					}).Return(&tc.previousRevision, nil)
					onlineStore.EXPECT().Purge(ctx, &tc.previousRevision).Return(nil)
				}
				metadataStore.EXPECT().UpdateFeatureGroup(ctx, types.UpdateFeatureGroupOpt{
					GroupName:        tc.group.Name,
					OnlineRevisionId: &tc.revision.ID,
				})
			}

			err := store.Materialize(ctx, tc.opt)
			if tc.expectedError != nil {
				assert.Error(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func prepareRevision(id int32, revision int64) types.Revision {
	return types.Revision{
		ID:        id,
		Revision:  revision,
		GroupName: "device_info",
		DataTable: fmt.Sprintf("device_info_%d", revision),
	}
}
func prepareGroup(revisionId *int32) types.FeatureGroup {
	return types.FeatureGroup{
		Name:             "device_info",
		OnlineRevisionID: revisionId,
		EntityName:       "device",
	}
}

func int64Ptr(i int64) *int64 {
	return &i
}

func int32Ptr(i int32) *int32 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
