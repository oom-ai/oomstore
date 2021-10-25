package redis

import (
	"github.com/go-redis/redis/v8"
)

type DB struct {
	redis.Client
}
