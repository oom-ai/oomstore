package redis_test

import (
	"context"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/internal/database/online/redis"
	"github.com/ethhte88/oomstore/internal/database/online/test_impl"
	"github.com/ethhte88/oomstore/internal/database/test/runtime_redis"
)

func prepareStore() (context.Context, online.Store) {
	ctx := context.Background()
	store := redis.Open(&runtime_redis.RedisDbOpt)
	if _, err := store.FlushDB(ctx).Result(); err != nil {
		panic(err)
	}

	return ctx, store
}

func TestOpen(t *testing.T) {
	test_impl.TestOpen(t, prepareStore)
}

func TestGetExisted(t *testing.T) {
	test_impl.TestGetExisted(t, prepareStore)
}

func TestGetNotExistedEntityKey(t *testing.T) {
	test_impl.TestGetNotExistedEntityKey(t, prepareStore)
}

func TestMultiGet(t *testing.T) {
	test_impl.TestMultiGet(t, prepareStore)
}

func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
	test_impl.TestPurgeRemovesSpecifiedRevision(t, prepareStore)
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
	test_impl.TestPurgeNotRemovesOtherRevisions(t, prepareStore)
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore)
}
