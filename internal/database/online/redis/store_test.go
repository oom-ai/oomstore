package redis_test

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/redis"
	"github.com/oom-ai/oomstore/internal/database/online/test_impl"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_redis"
)

func prepareStore(t *testing.T) (context.Context, online.Store) {
	ctx := context.Background()
	store := redis.Open(runtime_redis.GetOpt())
	if _, err := store.FlushDB(ctx).Result(); err != nil {
		t.Fatal(err)
	}

	return ctx, store
}

func destroyStore() {
	store := redis.Open(runtime_redis.GetOpt())
	if _, err := store.Client.FlushDB(context.Background()).Result(); err != nil {
		panic(err)
	}
}

func TestOpen(t *testing.T) {
	test_impl.TestOpen(t, prepareStore, destroyStore)
}

func TestGetExisted(t *testing.T) {
	test_impl.TestGetExisted(t, prepareStore, destroyStore)
}

func TestGetNotRevision(t *testing.T) {
	test_impl.TestGetNoRevision(t, prepareStore, destroyStore)
}

func TestGetNotExistedEntityKey(t *testing.T) {
	test_impl.TestGetNotExistedEntityKey(t, prepareStore, destroyStore)
}

func TestMultiGet(t *testing.T) {
	test_impl.TestMultiGet(t, prepareStore, destroyStore)
}

func TestGetByGroup(t *testing.T) {
	test_impl.TestGetByGroup(t, prepareStore, destroyStore)
}

func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
	test_impl.TestPurgeRemovesSpecifiedRevision(t, prepareStore, destroyStore)
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
	test_impl.TestPurgeNotRemovesOtherRevisions(t, prepareStore, destroyStore)
}

func TestPush(t *testing.T) {
	test_impl.TestPush(t, prepareStore, destroyStore)
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore, destroyStore)
}
