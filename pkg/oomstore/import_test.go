package oomstore_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/metadata"
	"github.com/ethhte88/oomstore/internal/database/metadata/mock_metadata"
	"github.com/ethhte88/oomstore/internal/database/offline/mock_offline"
	"github.com/ethhte88/oomstore/pkg/oomstore"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestChannelImportWithDependencyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.TEST__New(nil, offlineStore, metadataStore)

	testCases := []struct {
		description    string
		opt            types.ChannelImport
		mockFunc       func()
		wantRevisionID int
		wantError      error
	}{
		{
			description: "GetGroup failed",
			opt: types.ChannelImport{
				GroupName: "device_info",
			},
			mockFunc: func() {
				metadataStore.EXPECT().GetGroupByName(gomock.Any(), "device_info").Return(nil, fmt.Errorf("error"))
			},
			wantRevisionID: 0,
			wantError:      fmt.Errorf("error"),
		},
		{
			description: "ListFeature failed",
			opt: types.ChannelImport{
				GroupName: "device_info",
			},
			mockFunc: func() {
				metadataStore.EXPECT().GetGroupByName(gomock.Any(), "device_info").Return(&types.Group{ID: 1}, nil)
				metadataStore.EXPECT().ListFeature(gomock.Any(), metadata.ListFeatureOpt{GroupID: intPtr(1)}).Return(nil, nil)
			},
			wantRevisionID: 0,
			wantError:      fmt.Errorf("no features under group: device_info"),
		},
		{
			description: "GetEntity failed",
			opt: types.ChannelImport{
				GroupName: "device_info",
			},
			mockFunc: func() {
				metadataStore.EXPECT().GetGroupByName(gomock.Any(), "device_info").Return(&types.Group{ID: 1, EntityID: 1}, nil)
				metadataStore.EXPECT().ListFeature(gomock.Any(), gomock.Any()).Return(types.FeatureList{}, nil)
			},
			wantRevisionID: 0,
			wantError:      fmt.Errorf("no entity found by group: device_info"),
		},
		{
			description: "Create Revision failed",
			opt: types.ChannelImport{
				DataSource: types.CsvDataSource{
					Reader: strings.NewReader(`
device,model,price
1234,xiaomi,200
1235,apple,299
`),
					Delimiter: ",",
				},
				GroupName: "device_info",
			},
			mockFunc: func() {
				metadataStore.EXPECT().GetGroupByName(gomock.Any(), "device_info").Return(&types.Group{ID: 1, EntityID: 1, Entity: &types.Entity{Name: "device"}}, nil)
				metadataStore.EXPECT().ListFeature(gomock.Any(), metadata.ListFeatureOpt{GroupID: intPtr(1)}).
					Return(types.FeatureList{
						{
							Name: "model",
						},
						{
							Name: "price",
						},
					}, nil)
				offlineStore.EXPECT().Import(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), nil)

				metadataStore.EXPECT().CreateRevision(gomock.Any(), gomock.Any()).Return(0, "", fmt.Errorf("error"))
			},
			wantRevisionID: 0,
			wantError:      fmt.Errorf("error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().Refresh().Return(nil).AnyTimes()
			tc.mockFunc()
			revisionID, err := store.ChannelImport(context.Background(), tc.opt)
			assert.EqualError(t, err, tc.wantError.Error())
			assert.Equal(t, tc.wantRevisionID, revisionID)
		})
	}
}

func TestChannelImport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.TEST__New(nil, offlineStore, metadataStore)
	ctx := context.Background()

	testCases := []struct {
		description string

		opt        types.ChannelImport
		features   types.FeatureList
		entityID   int
		Entity     types.Entity
		header     []string
		revisionID int
		wantError  error
	}{
		{
			description: "import batch feature, succeed",
			opt: types.ChannelImport{
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
			opt: types.ChannelImport{
				GroupName: "device",
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
			opt: types.ChannelImport{
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
			metadataStore.EXPECT().Refresh().Return(nil).AnyTimes()
			metadataStore.EXPECT().GetGroupByName(ctx, tc.opt.GroupName).Return(&types.Group{
				ID:       1,
				Name:     tc.opt.GroupName,
				EntityID: tc.entityID,
				Entity:   &tc.Entity,
			}, nil)

			metadataStore.EXPECT().ListFeature(ctx, metadata.ListFeatureOpt{GroupID: intPtr(1)}).Return(tc.features, nil)

			metadataStore.EXPECT().CreateRevision(ctx, metadata.CreateRevisionOpt{
				Revision:    0,
				GroupID:     1,
				Description: tc.opt.Description,
			}).Return(tc.revisionID, "offline_1_1", nil).AnyTimes()

			offlineStore.EXPECT().Import(ctx, gomock.Any()).Return(int64(1000), nil).AnyTimes()

			if tc.opt.Revision == nil {
				metadataStore.EXPECT().UpdateRevision(ctx, metadata.UpdateRevisionOpt{
					RevisionID:  tc.revisionID,
					NewRevision: int64Ptr(1000),
				}).Return(nil).AnyTimes()
			}

			revisionID, err := store.ChannelImport(ctx, tc.opt)
			if tc.wantError != nil {
				assert.EqualError(t, err, tc.wantError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, revisionID, tc.revisionID)
		})
	}
}
