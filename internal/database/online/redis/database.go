package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/onestore-ai/onestore/internal/database/online"
)

var _ online.Store = &DB{}

const PipelineBatchSize = 10
const SeralizeIntBase = 36

type DB struct {
	*redis.Client
}

type RedisOpt struct {
	Host string
	Port int
	Pass string
	DB   int
}

func Open(opt *RedisOpt) *DB {
	redisOpt := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		Password: opt.Pass,
		DB:       opt.DB,
	}
	return &DB{redis.NewClient(&redisOpt)}
}
