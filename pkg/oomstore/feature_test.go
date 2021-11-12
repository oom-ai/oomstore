package oomstore_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	mock_metadatav2 "github.com/oom-ai/oomstore/internal/database/metadatav2/mock_metadata"
	"github.com/oom-ai/oomstore/internal/database/offline/mock_offline"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

func TestImportBatchFeatureWithDependencyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	offlineStore := mock_offline.NewMockStore(ctrl)
	metadatav2Store := mock_metadatav2.NewMockStore(ctrl)
	store := oomstore.NewOomStore(nil, offlineStore, metadatav2Store)

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
				metadatav2Store.EXPECT().
					ListFeature(gomock.Any(), metadatav2.ListFeatureOpt{GroupID: int16Ptr(1)}).
					Return(nil, fmt.Errorf("error"))
			},
			wantRevisionID: 0,
			wantError:      fmt.Errorf("error"),
		},
		{
			description: "GetFeatureGroup failed",
			opt:         types.ImportBatchFeaturesOpt{GroupID: 1},
			mockFunc: func() {
				metadatav2Store.EXPECT().
					ListFeature(gomock.Any(), metadatav2.ListFeatureOpt{GroupID: int16Ptr(1)}).
					Return(nil, nil)
				metadatav2Store.EXPECT().
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
				metadatav2Store.EXPECT().
					ListFeature(gomock.Any(), gomock.Any()).
					Return(nil, nil)
				metadatav2Store.EXPECT().
					GetFeatureGroup(gomock.Any(), 1).
					Return(&typesv2.FeatureGroup{ID: 1, EntityID: 1}, nil)
				metadatav2Store.EXPECT().
					GetEntity(gomock.Any(), 1).
					Return(nil, fmt.Errorf("error"))
			},
			wantRevisionID: 0,
			wantError:      fmt.Errorf("error"),
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
				metadatav2Store.EXPECT().
					ListFeature(gomock.Any(), metadatav2.ListFeatureOpt{GroupID: int16Ptr(1)}).
					Return(typesv2.FeatureList{
						{
							Name: "model",
						},
						{
							Name: "price",
						},
					}, nil)
				metadatav2Store.EXPECT().
					GetFeatureGroup(gomock.Any(), 1).
					Return(&typesv2.FeatureGroup{ID: 1, EntityID: 1}, nil)
				metadatav2Store.EXPECT().
					GetEntity(gomock.Any(), 1).
					Return(&typesv2.Entity{Name: "device"}, nil)

				offlineStore.
					EXPECT().
					Import(gomock.Any(), gomock.Any()).
					AnyTimes().Return(int64(1), "datatable", nil)

				metadatav2Store.EXPECT().
					CreateRevision(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("error"))
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
	metadatav2Store := mock_metadatav2.NewMockStore(ctrl)
	store := oomstore.NewOomStore(nil, offlineStore, metadatav2Store)

	testCases := []struct {
		description string

		opt              types.ImportBatchFeaturesOpt
		features         typesv2.FeatureList
		entityID         int16
		header           []string
		revisionID       int64
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
			features: typesv2.FeatureList{
				{
					Name: "model",
				},
				{
					Name: "price",
				},
			},
			entityID:         1,
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
			features: typesv2.FeatureList{
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
			description: "import batch feature: csv heaer of the data source doesn't match the feature group schema",
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
			features: typesv2.FeatureList{
				{
					Name: "model",
				},
			},
			entityID:         1,
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

			metadatav2Store.
				EXPECT().
				GetFeatureGroup(gomock.Any(), tc.opt.GroupID).
				Return(&typesv2.FeatureGroup{
					ID:       tc.opt.GroupID,
					EntityID: tc.entityID,
				}, nil)

			metadatav2Store.
				EXPECT().
				GetEntity(gomock.Any(), tc.entityID).
				Return(&typesv2.Entity{ID: tc.entityID}, nil)

			metadatav2Store.
				EXPECT().
				ListFeature(gomock.Any(), metadatav2.ListFeatureOpt{
					GroupID: &tc.opt.GroupID,
				}).
				Return(tc.features, nil)

			metadatav2Store.EXPECT().CreateRevision(gomock.Any(), metadatav2.CreateRevisionOpt{
				Revision:    int64(1),
				GroupID:     tc.opt.GroupID,
				DataTable:   stringPtr("datatable"),
				Description: tc.opt.Description,
			}).AnyTimes().Return(&typesv2.Revision{
				ID: int32(tc.revisionID),
			}, nil)

			revisionID, err := store.ImportBatchFeatures(context.Background(), tc.opt)

			assert.Equal(t, tc.wantError, err)
			assert.Equal(t, revisionID == 0, tc.revisionIDIsZero)
		})
	}
}

func TestCreateBatchFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	offlineStore := mock_offline.NewMockStore(ctrl)
	metadatav2Store := mock_metadatav2.NewMockStore(ctrl)
	store := oomstore.NewOomStore(nil, offlineStore, metadatav2Store)

	testCases := []struct {
		description string
		opt         metadatav2.CreateFeatureOpt
		valueType   string
		group       typesv2.FeatureGroup
		expectError bool
	}{
		{
			description: "create batch feature, succeed",
			opt: metadatav2.CreateFeatureOpt{
				Name:        "model",
				GroupID:     1,
				DBValueType: "VARCHAR(32)",
			},
			valueType: types.STRING,
			group: typesv2.FeatureGroup{
				Name:     "device_info",
				Category: types.BatchFeatureCategory,
			},
			expectError: false,
		},
		{
			description: "create stream feature, fail",
			opt: metadatav2.CreateFeatureOpt{
				Name:        "model",
				GroupID:     1,
				DBValueType: "BIGINT",
			},
			valueType: types.INT64,
			group: typesv2.FeatureGroup{
				Name:     "device_info",
				Category: types.StreamFeatureCategory,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			offlineStore.EXPECT().TypeTag(tc.opt.DBValueType).Return(tc.valueType, nil)
			metadatav2Store.EXPECT().GetFeatureGroup(gomock.Any(), tc.opt.GroupID).Return(&tc.group, nil)
			if tc.group.Category == types.BatchFeatureCategory {
				metadatav2Store.EXPECT().CreateFeature(gomock.Any(), tc.opt).Return(nil)
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
