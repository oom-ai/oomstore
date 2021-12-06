package oomstore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/metadata/mock_metadata"
	"github.com/ethhte88/oomstore/internal/database/offline/mock_offline"
	"github.com/ethhte88/oomstore/internal/database/online/mock_online"
	"github.com/ethhte88/oomstore/pkg/oomstore"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	prepare := func() (*oomstore.OomStore, *mock_online.MockStore, *mock_offline.MockStore, *mock_metadata.MockStore) {
		online := mock_online.NewMockStore(ctrl)
		offline := mock_offline.NewMockStore(ctrl)
		metadata := mock_metadata.NewMockStore(ctrl)
		oomstore := oomstore.TEST__New(online, offline, metadata)
		return oomstore, online, offline, metadata
	}

	t.Run("succeeded when all stores are available", func(t *testing.T) {
		store, online, offline, meta := prepare()
		online.EXPECT().Ping(ctx).Return(nil).AnyTimes()
		offline.EXPECT().Ping(ctx).Return(nil).AnyTimes()
		meta.EXPECT().Ping(ctx).Return(nil).AnyTimes()
		require.NoError(t, store.Ping(ctx))
	})

	t.Run("failed when online store is unavailable", func(t *testing.T) {
		store, online, offline, meta := prepare()
		online.EXPECT().Ping(ctx).Return(fmt.Errorf("whatever"))
		offline.EXPECT().Ping(ctx).Return(nil).AnyTimes()
		meta.EXPECT().Ping(ctx).Return(nil).AnyTimes()
		require.Error(t, store.Ping(ctx))
	})

	t.Run("failed when offline store is unavailable", func(t *testing.T) {
		store, online, offline, meta := prepare()
		online.EXPECT().Ping(ctx).Return(nil).AnyTimes()
		offline.EXPECT().Ping(ctx).Return(fmt.Errorf("whatever"))
		meta.EXPECT().Ping(ctx).Return(nil).AnyTimes()
		require.Error(t, store.Ping(ctx))
	})

	t.Run("failed when metadata store is unavailable", func(t *testing.T) {
		store, online, offline, meta := prepare()
		online.EXPECT().Ping(ctx).Return(nil).AnyTimes()
		offline.EXPECT().Ping(ctx).Return(nil).AnyTimes()
		meta.EXPECT().Ping(ctx).Return(fmt.Errorf("whatever"))
		require.Error(t, store.Ping(ctx))
	})
}
