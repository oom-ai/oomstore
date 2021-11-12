package oomstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	mock_metadatav2 "github.com/oom-ai/oomstore/internal/database/metadatav2/mock_metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/mock_offline"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/mock_online"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
	"github.com/stretchr/testify/assert"
)

func TestSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	onlineStore := mock_online.NewMockStore(ctrl)
	offlineStore := mock_offline.NewMockStore(ctrl)
	metadatav2Store := mock_metadatav2.NewMockStore(ctrl)
	store := oomstore.NewOomStore(onlineStore, offlineStore, metadatav2Store)
	ctx := context.Background()

	revision1 := prepareRevision(1, 10)
	revision2 := prepareRevision(2, 20)
	revision3 := prepareRevision(3, 30)
	entity := typesv2.Entity{
		ID:   1,
		Name: "device",
	}
	features := prepareFeatures(true, true)
	stream := make(chan *types.RawFeatureValueRecord)

	testCases := []struct {
		description      string
		opt              types.SyncOpt
		group            typesv2.FeatureGroup
		revision         typesv2.Revision
		previousRevision typesv2.Revision
		expectedError    error
	}{
		{
			description: "user-defined revision, succeed",
			opt: types.SyncOpt{
				GroupID:    1,
				RevisionID: 10,
			},
			group:            prepareGroup(int32Ptr(2)),
			revision:         revision1,
			previousRevision: revision2,
			expectedError:    nil,
		},
		{
			description: "latest revision, succeed",
			opt: types.SyncOpt{
				GroupID: 1,
			},
			group:            prepareGroup(int32Ptr(2)),
			revision:         revision3,
			previousRevision: revision2,
			expectedError:    nil,
		},
		{
			description: "no previous revision, succeed",
			opt: types.SyncOpt{
				GroupID: 1,
			},
			group:         prepareGroup(nil),
			revision:      revision1,
			expectedError: nil,
		},
		{
			description: "already in latest revision, fail",
			opt: types.SyncOpt{
				GroupID: 1,
			},
			group:         prepareGroup(int32Ptr(3)),
			revision:      revision3,
			expectedError: fmt.Errorf("online store already in the latest revision"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadatav2Store.EXPECT().GetFeatureGroup(ctx, tc.opt.GroupID).Return(&tc.group, nil)
			metadatav2Store.EXPECT().GetEntity(ctx, entity.ID).AnyTimes().Return(&entity, nil)
			metadatav2Store.EXPECT().ListFeature(ctx, metadatav2.ListFeatureOpt{GroupID: &tc.opt.GroupID}).AnyTimes().Return(features, nil)

			metadatav2Store.EXPECT().GetRevision(ctx, metadatav2.GetRevisionOpt{
				GroupID:    &tc.opt.GroupID,
				RevisionID: &tc.opt.RevisionID,
			}).Return(&tc.revision, nil)
			if tc.expectedError == nil {
				offlineStore.EXPECT().Export(ctx, offline.ExportOpt{
					DataTable:    tc.revision.DataTable,
					EntityName:   entity.Name,
					FeatureNames: features.Names(),
				}).Return(stream, nil)
				onlineStore.EXPECT().Import(ctx, online.ImportOpt{
					FeatureList: features,
					Revision:    &tc.revision,
					Entity:      &entity,
					Stream:      stream,
				})
				if tc.group.OnlineRevisionID != nil {
					metadatav2Store.EXPECT().GetRevision(ctx, metadatav2.GetRevisionOpt{
						RevisionID: tc.group.OnlineRevisionID,
					}).Return(&tc.previousRevision, nil)
					onlineStore.EXPECT().Purge(ctx, &tc.previousRevision).Return(nil)
				}
				metadatav2Store.EXPECT().UpdateFeatureGroup(ctx, metadatav2.UpdateFeatureGroupOpt{
					GroupID:             tc.group.ID,
					NewOnlineRevisionID: &tc.revision.ID,
				})
				metadatav2Store.EXPECT().
					UpdateRevision(gomock.Any(), gomock.Any()).
					AnyTimes().
					Return(int64(0), nil)
			}

			err := store.Sync(ctx, tc.opt)
			if tc.expectedError != nil {
				assert.Error(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func prepareRevision(id int32, revision int64) typesv2.Revision {
	return typesv2.Revision{
		ID:        id,
		Revision:  revision,
		GroupID:   1,
		DataTable: fmt.Sprintf("device_info_%d", revision),
	}
}
func prepareGroup(revisionId *int32) typesv2.FeatureGroup {
	return typesv2.FeatureGroup{
		Name:             "device_info",
		OnlineRevisionID: revisionId,
		EntityID:         1,
	}
}

func int16Ptr(i int16) *int16 {
	return &i
}

func int32Ptr(i int32) *int32 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
