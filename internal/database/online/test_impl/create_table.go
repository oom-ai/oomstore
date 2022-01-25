package test_impl

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/stretchr/testify/assert"
)

func TestPrepareStreamTable(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	for _, group := range simpleStreamData.groups {
		t.Run("create stream table", func(t *testing.T) {
			err := store.PrepareStreamTable(ctx, online.PrepareStreamTableOpt{
				EntityName: group.Entity.Name,
				GroupID:    group.ID,
			})
			assert.NoError(t, err, "create stream table failed: %v", err)
		})
	}

	for _, feature := range simpleStreamData.features {
		t.Run("stream table add column", func(t *testing.T) {
			err := store.PrepareStreamTable(ctx, online.PrepareStreamTableOpt{
				EntityName: feature.Entity().Name,
				GroupID:    feature.GroupID,
				Feature:    feature,
			})
			assert.NoError(t, err, "stream table add column failed: %v", err)
		})
	}
}
