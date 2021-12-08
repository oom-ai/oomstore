package test_impl

import (
	"context"
	"encoding/csv"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

type PrepareStoreFn func() (context.Context, offline.Store)

func buildTestDataTable(ctx context.Context, t *testing.T, store offline.Store, features []*types.Feature, dataTable string, reader *csv.Reader) {
	entity := &types.Entity{
		Name:   "device",
		Length: 10,
	}
	header := []string{"device"}
	for _, f := range features {
		header = append(header, f.Name)
	}
	opt := offline.ImportOpt{
		Entity:        entity,
		DataTableName: dataTable,
		Features:      features,
		Header:        header,
		CsvReader:     reader,
	}
	_, err := store.Import(ctx, opt)
	require.NoError(t, err)
}
