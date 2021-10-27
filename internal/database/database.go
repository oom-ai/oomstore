package database

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	metadataPG "github.com/oom-ai/oomstore/internal/database/metadata/postgres"

	"github.com/oom-ai/oomstore/internal/database/offline"
	offlinePG "github.com/oom-ai/oomstore/internal/database/offline/postgres"

	"github.com/oom-ai/oomstore/internal/database/online"
	onlinePG "github.com/oom-ai/oomstore/internal/database/online/postgres"
	onlineRedis "github.com/oom-ai/oomstore/internal/database/online/redis"
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
		return metadataPG.Open(opt.PostgresDbOpt)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func CreateMetadataDatabase(ctx context.Context, opt types.MetaStoreOpt) error {
	switch opt.Backend {
	case types.POSTGRES:
		return metadataPG.CreateDatabase(ctx, *opt.PostgresDbOpt)
	default:
		return fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func OpenOfflineStore(opt types.OfflineStoreOpt) (offline.Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return offlinePG.Open(opt.PostgresDbOpt)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}
