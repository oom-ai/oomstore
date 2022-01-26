package test_impl

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/stretchr/testify/assert"
)

func TestCreateTable(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	t.Run("create stream table", func(t *testing.T) {
		err := store.CreateTable(ctx, online.CreateTableOpt{
			EntityName: SampleStream.Entity.Name,
			TableName:  "stream_online",
			Features:   SampleStream.Features,
		})
		assert.NoError(t, err, "create stream table failed: %v", err)
	})
	t.Run("create batch table", func(t *testing.T) {
		err := store.CreateTable(ctx, online.CreateTableOpt{
			EntityName: SampleSmall.Entity.Name,
			TableName:  "batch_online",
			Features:   SampleSmall.Features,
		})
		assert.NoError(t, err, "create batch table failed: %v", err)
	})
}
