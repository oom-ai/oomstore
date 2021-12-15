package redis

import (
	"context"
	"fmt"

	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/go-redis/redis/v8"
)

const PipelineBatchSize = 10
const SerializeIntBase = 36

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
