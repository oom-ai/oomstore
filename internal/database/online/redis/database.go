package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

const PipelineBatchSize = 10
const SeralizeIntBase = 36

type DB struct {
	*redis.Client
}

func Open(opt *types.RedisDbOpt) *DB {
	redisOpt := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", opt.Host, opt.Port),
		Password: opt.Pass,
		DB:       opt.Database,
	}
	return &DB{redis.NewClient(&redisOpt)}
}
