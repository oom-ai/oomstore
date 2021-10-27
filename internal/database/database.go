package database

import (
	"context"
	"fmt"

	"github.com/onestore-ai/onestore/pkg/onestore/types"

	"github.com/onestore-ai/onestore/internal/database/metadata"
	metadataPg "github.com/onestore-ai/onestore/internal/database/metadata/postgres"

	"github.com/onestore-ai/onestore/internal/database/online"
	onlinePG "github.com/onestore-ai/onestore/internal/database/online/postgres"
	onlineRedis "github.com/onestore-ai/onestore/internal/database/online/redis"
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

func OpenMetadataStore(opt types.MetaStoreOpt) (metadata.Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return metadataPg.Open(opt.PostgresDbOpt)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func CreateMetadataDatabase(ctx context.Context, opt types.MetaStoreOpt) error {
	switch opt.Backend {
	case types.POSTGRES:
		return metadataPg.CreateDatabase(ctx, *opt.PostgresDbOpt)
	default:
		return fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}
