package test_impl

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

type PrepareStoreFn func(t *testing.T) (context.Context, offline.Store)

func buildTestDataTable(ctx context.Context, t *testing.T, store offline.Store, features []*types.Feature, dataTable string, source *offline.CSVSource) {
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
		Source:        source,
	}
	_, err := store.Import(ctx, opt)
	require.NoError(t, err)
}
