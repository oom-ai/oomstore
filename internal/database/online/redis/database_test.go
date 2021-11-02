package redis

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/test"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_redis"
)

func prepareStore() (context.Context, online.Store) {
	ctx := context.Background()
	store := Open(&runtime_redis.RedisDbOpt)
	if _, err := store.FlushDB(ctx).Result(); err != nil {
		panic(err)
	}

	return ctx, store
}

func TestOpen(t *testing.T) {
	test.TestOpen(t, prepareStore)
}

func TestGetExisted(t *testing.T) {
	test.TestGetExisted(t, prepareStore)
}

func TestGetNotExistedEntityKey(t *testing.T) {
	test.TestGetNotExistedEntityKey(t, prepareStore)
}

func TestMultiGet(t *testing.T) {
	test.TestMultiGet(t, prepareStore)
}

func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
	test.TestPurgeRemovesSpecifiedRevision(t, prepareStore)
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
	test.TestPurgeNotRemovesOtherRevisions(t, prepareStore)
}
