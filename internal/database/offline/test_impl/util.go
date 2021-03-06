package test_impl

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

type PrepareStoreFn func(t *testing.T) (context.Context, offline.Store)

type DestroyStoreFn func()

func buildTestSnapshotTable(
	ctx context.Context,
	t *testing.T,
	store offline.Store,
	features []*types.Feature,
	revision int64,
	snapshotTable string,
	source *offline.CSVSource,
) {
	entity := &types.Entity{Name: "device"}
	header := []string{"device"}
	for _, f := range features {
		header = append(header, f.Name)
	}
	opt := offline.ImportOpt{
		EntityName:        entity.Name,
		SnapshotTableName: snapshotTable,
		Features:          features,
		Header:            header,
		Source:            source,
		Revision:          &revision,
		Category:          types.CategoryStream,
	}
	_, err := store.Import(ctx, opt)
	require.NoError(t, err)
}
