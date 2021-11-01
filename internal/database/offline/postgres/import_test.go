package postgres_test

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestImport(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	entity := types.Entity{
		Name:   "device",
		Length: 16,
	}

	tempFeaturesFile := tempFeatureFile(t)
	opt := offline.ImportOpt{
		ImportBatchFeaturesOpt: types.ImportBatchFeaturesOpt{
			GroupName:   "device",
			Description: "description",
			DataSource: types.LocalFileDataSourceOpt{
				FilePath:  tempFeaturesFile,
				Delimiter: ",",
			},
		},
		Entity: &entity,
		Features: []*types.Feature{
			{
				Name:        "model",
				DBValueType: "invalid-db-value-type"},
			{
				Name:        "price",
				DBValueType: "int",
			},
		},
		Header: []string{"device", "model", "price"},
	}

	t.Run("invalid db value type", func(t *testing.T) {
		_, _, err := db.Import(context.Background(), opt)
		assert.NotNil(t, err)
	})

	t.Run("normal import call", func(t *testing.T) {
		opt.Features[0].DBValueType = "varchar(32)"
		_, tableName, err := db.Import(context.Background(), opt)
		assert.Nil(t, err)

		type T struct {
			Device string `db:"device"`
			Model  string `db:"model"`
			Price  int    `db:"price"`
		}
		records := make([]T, 0)

		assert.Nil(t, db.SelectContext(context.Background(), &records, fmt.Sprintf("select * from %s", tableName)))
		assert.Equal(t, 4, len(records))

		sort.Slice(records, func(i, j int) bool {
			return records[i].Device < records[j].Device
		})
		assert.Equal(t, T{Device: "1234", Model: "xiaomi", Price: 1899}, records[0])
		assert.Equal(t, T{Device: "1235", Model: "apple", Price: 4999}, records[1])
		assert.Equal(t, T{Device: "1236", Model: "huawei", Price: 5999}, records[2])
		assert.Equal(t, T{Device: "1237", Model: "oneplus", Price: 3999}, records[3])
	})
}

func tempFeatureFile(t *testing.T) string {
	filePath := fmt.Sprintf("%s/features.csv", t.TempDir())
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	records := [][]string{
		{"device", "model", "price"},
		{"1234", "xiaomi", "1899"},
		{"1235", "apple", "4999"},
		{"1236", "huawei", "5999"},
		{"1237", "oneplus", "3999"},
	}

	if err = w.WriteAll(records); err != nil {
		panic(err)
	}
	w.Flush()
	return file.Name()
}
