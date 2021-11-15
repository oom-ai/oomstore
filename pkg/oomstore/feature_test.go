package oomstore_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/mock_metadata"
	"github.com/oom-ai/oomstore/internal/database/offline/mock_offline"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
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
		wantRevisionID int32
		wantError      error
	}{
		{
			description: "ListFeature failed",
			opt:         types.ImportBatchFeaturesOpt{GroupID: 1},
			mockFunc: func() {
				metadataStore.EXPECT().
					ListFeature(gomock.Any(), metadata.ListFeatureOpt{GroupID: int16Ptr(1)}).
					Return(nil)
			},
			wantRevisionID: 0,
			wantError:      fmt.Errorf("no featues under group id: '1'"),
		},
		{
			description: "GetFeatureGroup failed",
			opt:         types.ImportBatchFeaturesOpt{GroupID: 1},
			mockFunc: func() {
				metadataStore.EXPECT().
					ListFeature(gomock.Any(), metadata.ListFeatureOpt{GroupID: int16Ptr(1)}).
					Return(types.FeatureList{})
				metadataStore.EXPECT().
					GetFeatureGroup(gomock.Any(), int16(1)).
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
					GetFeatureGroup(gomock.Any(), int16(1)).
					Return(&types.FeatureGroup{ID: 1, EntityID: 1}, nil)
				metadataStore.EXPECT().
					GetEntity(gomock.Any(), int16(1)).
					Return(nil, fmt.Errorf("error"))
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
					ListFeature(gomock.Any(), metadata.ListFeatureOpt{GroupID: int16Ptr(1)}).
					Return(types.FeatureList{
						{
							Name: "model",
						},
						{
							Name: "price",
						},
					})
				metadataStore.EXPECT().
					GetFeatureGroup(gomock.Any(), int16(1)).
					Return(&types.FeatureGroup{ID: 1, EntityID: 1, Entity: &types.Entity{Name: "device"}}, nil)
				metadataStore.EXPECT().
					GetEntity(gomock.Any(), int16(1)).
					Return(&types.Entity{Name: "device"}, nil)
				offlineStore.
					EXPECT().
					Import(gomock.Any(), gomock.Any()).
					AnyTimes().Return(int64(1), "datatable", nil)

				metadataStore.EXPECT().
					CreateRevision(gomock.Any(), gomock.Any()).
					Return(int32(0), fmt.Errorf("error"))
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

	testCases := []struct {
		description string

		opt              types.ImportBatchFeaturesOpt
		features         types.FeatureList
		entityID         int16
		Entity           types.Entity
		header           []string
		revisionID       int32
		revisionIDIsZero bool
		wantError        error
	}{
		{
			description: "import batch feature: succeed",
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
			entityID:         1,
			Entity:           types.Entity{Name: "device"},
			header:           []string{"device", "model", "price"},
			revisionID:       1,
			revisionIDIsZero: false,
			wantError:        nil,
		},
		{
			description: "import batch feature: csv data source has duplicated columns",
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
			entityID:         1,
			header:           []string{"device", "model"},
			revisionID:       0,
			revisionIDIsZero: true,
			wantError:        fmt.Errorf("csv data source has duplicated columns: %v", []string{"device", "model", "model"}),
		},
		{
			description: "import batch feature: csv header of the data source doesn't match the feature group schema",
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
			entityID:         1,
			Entity:           types.Entity{Name: "device"},
			header:           []string{"device", "model", "price"},
			revisionID:       0,
			revisionIDIsZero: true,
			wantError: fmt.Errorf("csv header of the data source %v doesn't match the feature group schema %v",
				[]string{"device", "model", "price"},
				[]string{"device", "model"},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			offlineStore.
				EXPECT().
				Import(gomock.Any(), gomock.Any()).
				AnyTimes().Return(int64(1), "datatable", nil)

			metadataStore.
				EXPECT().
				GetFeatureGroup(gomock.Any(), tc.opt.GroupID).
				Return(&types.FeatureGroup{
					ID:       tc.opt.GroupID,
					EntityID: tc.entityID,
					Entity:   &tc.Entity,
				}, nil)

			metadataStore.
				EXPECT().
				GetEntity(gomock.Any(), tc.entityID).
				Return(&types.Entity{ID: tc.entityID, Name: "device"}, nil)
			metadataStore.
				EXPECT().
				ListFeature(gomock.Any(), metadata.ListFeatureOpt{
					GroupID: &tc.opt.GroupID,
				}).
				Return(tc.features)

			metadataStore.EXPECT().CreateRevision(gomock.Any(), metadata.CreateRevisionOpt{
				Revision:    int64(1),
				GroupID:     tc.opt.GroupID,
				DataTable:   stringPtr("datatable"),
				Description: tc.opt.Description,
			}).AnyTimes().
				Return(tc.revisionID, nil)

			revisionID, err := store.ImportBatchFeatures(context.Background(), tc.opt)
			if tc.wantError != nil {
				assert.EqualError(t, err, tc.wantError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, revisionID == 0, tc.revisionIDIsZero)
		})
	}
}

func TestCreateBatchFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.NewOomStore(nil, offlineStore, metadataStore)

	testCases := []struct {
		description string
		opt         metadata.CreateFeatureOpt
		valueType   string
		group       types.FeatureGroup
		expectError bool
	}{
		{
			description: "create batch feature, succeed",
			opt: metadata.CreateFeatureOpt{
				Name:        "model",
				GroupID:     1,
				DBValueType: "VARCHAR(32)",
				ValueType:   types.STRING,
			},
			valueType: types.STRING,
			group: types.FeatureGroup{
				Name:     "device_info",
				Category: types.BatchFeatureCategory,
			},
			expectError: false,
		},
		{
			description: "create stream feature, fail",
			opt: metadata.CreateFeatureOpt{
				Name:        "model",
				GroupID:     1,
				DBValueType: "BIGINT",
				ValueType:   types.INT64,
			},
			valueType: types.INT64,
			group: types.FeatureGroup{
				Name:     "device_info",
				Category: types.StreamFeatureCategory,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			metadataStore.EXPECT().
				GetFeatureGroup(gomock.Any(), tc.opt.GroupID).
				Return(&tc.group, nil)

			if tc.group.Category == types.BatchFeatureCategory {
				tc.opt.ValueType = tc.valueType

				offlineStore.EXPECT().
					TypeTag(tc.opt.DBValueType).
					Return(tc.valueType, nil)

				metadataStore.EXPECT().
					CreateFeature(gomock.Any(), tc.opt).
					Return(int16(0), nil)
			}

			_, err := store.CreateBatchFeature(context.Background(), tc.opt)
			if tc.expectError {
				assert.Error(t, err, fmt.Errorf("expected batch feature group, got %s feature group", tc.group.Category))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
