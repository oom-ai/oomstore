package database

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	metadataPG "github.com/oom-ai/oomstore/internal/database/metadata/postgres"
	metadatav2PG "github.com/oom-ai/oomstore/internal/database/metadata/postgres"

	"github.com/oom-ai/oomstore/internal/database/offline"
	offlinePG "github.com/oom-ai/oomstore/internal/database/offline/postgres"

	"github.com/oom-ai/oomstore/internal/database/online"
	onlinePG "github.com/oom-ai/oomstore/internal/database/online/postgres"
	onlineRedis "github.com/oom-ai/oomstore/internal/database/online/redis"
)

func OpenOnlineStore(opt types.OnlineStoreConfig) (online.Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return onlinePG.Open(opt.Postgres)
	case types.REDIS:
		return onlineRedis.Open(opt.Redis), nil
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func OpenMetadataStore(opt types.MetadataStoreConfig) (metadata.Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return metadataPG.Open(opt.Postgres)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func OpenMetadatav2Store(opt types.MetadataStoreConfig) (metadata.Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return metadatav2PG.Open(context.Background(), opt.Postgres)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func CreateMetadataDatabase(ctx context.Context, opt types.MetadataStoreConfig) error {
	switch opt.Backend {
	case types.POSTGRES:
		return metadatav2PG.CreateDatabase(ctx, *opt.Postgres)
	default:
		return fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func OpenOfflineStore(opt types.OfflineStoreConfig) (offline.Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return offlinePG.Open(opt.Postgres)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}
