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

func TestImportBatchFeatures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	offlineStore := mock_offline.NewMockStore(ctrl)
	metadataStore := mock_metadata.NewMockStore(ctrl)
	store := oomstore.NewOomStore(nil, offlineStore, metadataStore)

	testCases := []struct {
		description string

		opt        types.ImportBatchFeaturesOpt
		features   types.FeatureList
		entityName string
		header     []string
		wantError  error
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
			entityName: "device",
			header:     []string{"device", "model", "price"},
			wantError:  nil,
		},
		{
			description: "import batch feature: csv data source has duplicated columns",
			opt: types.ImportBatchFeaturesOpt{
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
			entityName: "device",
			header:     []string{"device", "model"},
			wantError:  fmt.Errorf("csv data source has duplicated columns: %v", []string{"device", "model", "model"}),
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
			features: types.FeatureList{
				{
					Name: "model",
				},
			},
			entityName: "device",
			header:     []string{"device", "model", "price"},
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
				GetFeatureGroup(gomock.Any(), tc.opt.GroupName).
				Return(&types.FeatureGroup{
					Name:       tc.opt.GroupName,
					EntityName: tc.entityName,
				}, nil)

			metadataStore.
				EXPECT().
				GetEntity(gomock.Any(), tc.entityName).
				Return(&types.Entity{Name: tc.entityName}, nil)

			metadataStore.
				EXPECT().
				ListFeature(gomock.Any(), types.ListFeatureOpt{
					GroupName: &tc.opt.GroupName,
				}).
				Return(tc.features, nil)

			metadataStore.EXPECT().CreateRevision(gomock.Any(), metadata.CreateRevisionOpt{
				Revision:    int64(1),
				GroupName:   tc.opt.GroupName,
				DataTable:   "datatable",
				Description: tc.opt.Description,
			}).AnyTimes().Return(nil)

			assert.Equal(t, tc.wantError, store.ImportBatchFeatures(context.Background(), tc.opt))
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
		opt         types.CreateFeatureOpt
		valueType   string
		group       types.FeatureGroup
		expectError bool
	}{
		{
			description: "create batch feature, succeed",
			opt: types.CreateFeatureOpt{
				FeatureName: "model",
				GroupName:   "device_info",
				DBValueType: "VARCHAR(32)",
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
			opt: types.CreateFeatureOpt{
				FeatureName: "model",
				GroupName:   "device_info",
				DBValueType: "BIGINT",
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
			offlineStore.EXPECT().TypeTag(tc.opt.DBValueType).Return(tc.valueType, nil)
			metadataStore.EXPECT().GetFeatureGroup(gomock.Any(), tc.opt.GroupName).Return(&tc.group, nil)
			if tc.group.Category == types.BatchFeatureCategory {
				metadataStore.EXPECT().CreateFeature(gomock.Any(), metadata.CreateFeatureOpt{
					CreateFeatureOpt: tc.opt,
					ValueType:        tc.valueType,
				}).Return(nil)
			}

			err := store.CreateBatchFeature(context.Background(), tc.opt)
			if tc.expectError {
				assert.Error(t, err, fmt.Errorf("expected batch feature group, got %s feature group", tc.group.Category))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
