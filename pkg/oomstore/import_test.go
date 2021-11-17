package oomstore_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/mock_metadata"
	"github.com/oom-ai/oomstore/internal/database/offline/mock_offline"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

func TestImportBatchFeatureWithDependencyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.NewOomStore(nil, offlineStore, metadataStore)

	testCases := []struct {
		description    string
		opt            types.ImportBatchFeaturesOpt
		mockFunc       func()
		wantRevisionID int
		wantError      error
	}{
		{
			description: "ListFeature failed",
			opt:         types.ImportBatchFeaturesOpt{GroupID: 1},
			mockFunc: func() {
				metadataStore.EXPECT().
					ListFeature(gomock.Any(), metadata.ListFeatureOpt{GroupID: intPtr(1)}).
					Return(nil)
			},
			wantRevisionID: 0,
			wantError:      fmt.Errorf("no features under group id: '1'"),
		},
		{
			description: "GetFeatureGroup failed",
			opt:         types.ImportBatchFeaturesOpt{GroupID: 1},
			mockFunc: func() {
				metadataStore.EXPECT().
					ListFeature(gomock.Any(), metadata.ListFeatureOpt{GroupID: intPtr(1)}).
					Return(types.FeatureList{})
				metadataStore.EXPECT().
					GetFeatureGroup(gomock.Any(), 1).
					Return(nil, fmt.Errorf("error"))
			},
			wantRevisionID: 0,
			wantError:      fmt.Errorf("error"),
		},
		{
			description: "GetEntity failed",
			opt:         types.ImportBatchFeaturesOpt{GroupID: 1},
			mockFunc: func() {
				metadataStore.EXPECT().
					ListFeature(gomock.Any(), gomock.Any()).
					Return(types.FeatureList{})
				metadataStore.EXPECT().
					GetFeatureGroup(gomock.Any(), 1).
					Return(&types.FeatureGroup{ID: 1, EntityID: 1}, nil)
			},
			wantRevisionID: 0,
			wantError:      fmt.Errorf("no entity found by group id: '1'"),
		},
		{
			description: "Create Revision failed",
			opt: types.ImportBatchFeaturesOpt{
				DataSource: types.CsvDataSource{
					Reader: strings.NewReader(`
device,model,price
1234,xiaomi,200
1235,apple,299
`),
					Delimiter: ",",
				},
				GroupID: 1,
			},
			mockFunc: func() {
				metadataStore.EXPECT().
					ListFeature(gomock.Any(), metadata.ListFeatureOpt{GroupID: intPtr(1)}).
					Return(types.FeatureList{
						{
							Name: "model",
						},
						{
							Name: "price",
						},
					})
				metadataStore.EXPECT().
					GetFeatureGroup(gomock.Any(), 1).
					Return(&types.FeatureGroup{ID: 1, EntityID: 1, Entity: &types.Entity{Name: "device"}}, nil)
				offlineStore.
					EXPECT().
					Import(gomock.Any(), gomock.Any()).
					AnyTimes().Return(int64(1), nil)

				metadataStore.EXPECT().
					CreateRevision(gomock.Any(), gomock.Any()).
					Return(0, "", fmt.Errorf("error"))
			},
			wantRevisionID: 0,
			wantError:      fmt.Errorf("error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc.mockFunc()
			revisionID, err := store.ImportBatchFeatures(context.Background(), tc.opt)
			assert.Equal(t, tc.wantError, err)
			assert.Equal(t, tc.wantRevisionID, revisionID)
		})
	}
}

func TestImportBatchFeatures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.NewOomStore(nil, offlineStore, metadataStore)
	ctx := context.Background()

	testCases := []struct {
		description string

		opt        types.ImportBatchFeaturesOpt
		features   types.FeatureList
		entityID   int
		Entity     types.Entity
		header     []string
		revisionID int
		wantError  error
	}{
		{
			description: "import batch feature, succeed",
			opt: types.ImportBatchFeaturesOpt{
				DataSource: types.CsvDataSource{
					Reader: strings.NewReader(`device,model,price
1234,xiaomi,200
1235,apple,299
`),
					Delimiter: ",",
				},
			},
			features: types.FeatureList{
				{
					Name: "model",
				},
				{
					Name: "price",
				},
			},
			entityID:   1,
			Entity:     types.Entity{Name: "device"},
			header:     []string{"device", "model", "price"},
			revisionID: 1,
			wantError:  nil,
		},
		{
			description: "import batch feature, csv data source has duplicated columns",
			opt: types.ImportBatchFeaturesOpt{
				GroupID: 1,
				DataSource: types.CsvDataSource{
					Reader: strings.NewReader(`
device,model,model
1234,xiaomi,xiaomi
1235,apple,xiaomi
`),
					Delimiter: ",",
				},
			},
			features: types.FeatureList{
				{
					Name: "model",
				},
				{
					Name: "price",
				},
			},
			entityID:   1,
			header:     []string{"device", "model"},
			revisionID: 0,
			wantError:  fmt.Errorf("csv data source has duplicated columns: %v", []string{"device", "model", "model"}),
		},
		{
			description: "import batch feature, csv header of the data source doesn't match the feature group schema",
			opt: types.ImportBatchFeaturesOpt{
				DataSource: types.CsvDataSource{
					Reader: strings.NewReader(`
device,model,price
1234,xiaomi,200
1235,apple,299
`),
					Delimiter: ",",
				},
			},
			features: types.FeatureList{
				{
					Name: "model",
				},
			},
			entityID:   1,
			Entity:     types.Entity{Name: "device"},
			header:     []string{"device", "model", "price"},
			revisionID: 0,
			wantError: fmt.Errorf("csv header of the data source %v doesn't match the feature group schema %v",
				[]string{"device", "model", "price"},
				[]string{"device", "model"},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().GetFeatureGroup(ctx, tc.opt.GroupID).
				Return(&types.FeatureGroup{
					ID:       tc.opt.GroupID,
					EntityID: tc.entityID,
					Entity:   &tc.Entity,
				}, nil)

			metadataStore.EXPECT().ListFeature(ctx, metadata.ListFeatureOpt{GroupID: &tc.opt.GroupID}).Return(tc.features)

			metadataStore.EXPECT().CreateRevision(ctx, metadata.CreateRevisionOpt{
				Revision:    0,
				GroupID:     tc.opt.GroupID,
				Description: tc.opt.Description,
			}).Return(tc.revisionID, "data_1_1", nil).AnyTimes()

			offlineStore.EXPECT().Import(ctx, gomock.Any()).Return(int64(1000), nil).AnyTimes()

			if tc.opt.Revision == nil {
				metadataStore.EXPECT().UpdateRevision(ctx, metadata.UpdateRevisionOpt{
					RevisionID:  tc.revisionID,
					NewRevision: int64Ptr(1000),
				}).Return(nil).AnyTimes()
			}

			revisionID, err := store.ImportBatchFeatures(ctx, tc.opt)
			if tc.wantError != nil {
				assert.EqualError(t, err, tc.wantError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, revisionID, tc.revisionID)
		})
	}
}
