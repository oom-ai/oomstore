package test_impl

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func TestPush(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	entity := types.Entity{
		ID:   18,
		Name: "user",
	}
	group := types.Group{
		ID:       1,
		Name:     "vip",
		Category: types.CategoryStream,
		EntityID: entity.ID,
	}
	features := types.FeatureList{
		&types.Feature{
			ID:        1,
			Name:      "age",
			ValueType: types.Int64,
			GroupID:   group.ID,
		},
		&types.Feature{
			ID:        2,
			Name:      "gender",
			ValueType: types.String,
			GroupID:   group.ID,
		},
	}
	revision := int64(2)
	records := []types.StreamRecord{
		{

			EntityKey: "1000",
			UnixMilli: 8000,
			Values:    []interface{}{18, "F"},
		},
		{

			EntityKey: "1001",
			UnixMilli: 8001,
			Values:    []interface{}{21, "M"},
		},
	}

	pushOpt := offline.PushOpt{
		GroupID:      group.ID,
		Revision:     revision,
		EntityName:   entity.Name,
		FeatureNames: features.Names(),
		Records:      records,
	}

	t.Run("push when cdc table not exists", func(t *testing.T) {
		err := store.Push(ctx, pushOpt)
		require.Error(t, err)
		require.True(t, errdefs.IsNotFound(err))
	})

	t.Run("push when cdc table exists", func(t *testing.T) {
		err := store.CreateTable(ctx, offline.CreateTableOpt{
			TableName: dbutil.OfflineStreamCdcTableName(group.ID, revision),
			Entity:    &entity,
			Features:  features,
			TableType: types.TableStreamCdc,
		})
		require.NoError(t, err)

		err = store.Push(ctx, pushOpt)
		require.NoError(t, err)
	})
}
