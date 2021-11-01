package redis

import (
	"strconv"
	"testing"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
	ctx, store := prepare()
	defer store.Close()
	revision := types.Revision{ID: 382}

	for i := 0; i < PipelineBatchSize+1; i++ {
		key, err := SerializeRedisKey(revision.ID, i)
		require.NoError(t, err)
		require.NoError(t, store.HSet(ctx, key, strconv.Itoa(i), strconv.Itoa(i+1)).Err())

		v, err := store.HGet(ctx, key, strconv.Itoa(i)).Result()
		require.NoError(t, err)
		require.Equal(t, v, strconv.Itoa(i+1))
	}

	err := store.Purge(ctx, &revision)
	require.NoError(t, err)

	sz, err := store.DBSize(ctx).Result()
	require.NoError(t, err)
	require.Equal(t, int64(0), sz)
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
	ctx, store := prepare()
	defer store.Close()
	revision := types.Revision{ID: 382}

	// prepare the revision to be purged
	for i := 0; i < PipelineBatchSize+1; i++ {
		key, err := SerializeRedisKey(revision.ID, i)
		require.NoError(t, err)
		require.NoError(t, store.HSet(ctx, key, strconv.Itoa(i), strconv.Itoa(i+1)).Err())

		v, err := store.HGet(ctx, key, strconv.Itoa(i)).Result()
		require.NoError(t, err)
		require.Equal(t, v, strconv.Itoa(i+1))
	}

	// prepare another revision
	for i := 0; i < 10; i++ {
		key, err := SerializeRedisKey(0, i)
		require.NoError(t, err)
		require.NoError(t, store.HSet(ctx, key, strconv.Itoa(i), strconv.Itoa(i+1)).Err())

		v, err := store.HGet(ctx, key, strconv.Itoa(i)).Result()
		require.NoError(t, err)
		require.Equal(t, v, strconv.Itoa(i+1))
	}

	err := store.Purge(ctx, &revision)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		key, err := SerializeRedisKey(0, i)
		require.NoError(t, err)
		v, err := store.HGet(ctx, key, strconv.Itoa(i)).Result()
		require.NoError(t, err)
		require.Equal(t, v, strconv.Itoa(i+1))
	}
}
