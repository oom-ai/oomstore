package redis

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/test/test_redis"
	"github.com/stretchr/testify/assert"
)

func prepare() (context.Context, *DB) {
	ctx := context.Background()
	store := Open(&test_redis.RedisDbOpt)
	if _, err := store.FlushDB(ctx).Result(); err != nil {
		panic(err)
	}
	return ctx, store
}

func TestOpen(t *testing.T) {
	ctx, store := prepare()
	res, err := store.Ping(ctx).Result()
	assert.Nil(t, err)
	assert.Equal(t, res, "PONG")
}
