package database

import (
	"fmt"

	"github.com/onestore-ai/onestore/internal/database/online"
	onlinePG "github.com/onestore-ai/onestore/internal/database/online/postgres"
	onlineRedis "github.com/onestore-ai/onestore/internal/database/online/redis"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func OpenOnlineStore(opt types.OnlineStoreOpt) (online.Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return onlinePG.Open(opt.PostgresDbOpt)
	case types.REDIS:
		return onlineRedis.Open(opt.RedisDbOpt), nil
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}
