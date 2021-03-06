package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	Backend           = types.BackendRedis
	PipelineBatchSize = 10
)

var _ online.Store = &DB{}

type DB struct {
	*redis.Client
}

func (db *DB) Ping(ctx context.Context) error {
	_, err := db.Client.Ping(ctx).Result()
	return err
}

func Open(opt *types.RedisOpt) *DB {
	redisOpt := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", opt.Host, opt.Port),
		Password: opt.Password,
		DB:       opt.Database,
	}
	return &DB{redis.NewClient(&redisOpt)}
}

func (db *DB) CreateTable(ctx context.Context, opt online.CreateTableOpt) error {
	return nil
}
