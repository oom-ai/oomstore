package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/onestore/types"
)

const PipelineBatchSize = 10
const SeralizeIntBase = 36

var _ online.Store = &DB{}

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
