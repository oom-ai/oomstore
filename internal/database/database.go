package database

import (
	"context"
	"fmt"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"

	"github.com/ethhte88/oomstore/internal/database/metadata"
	metadataPG "github.com/ethhte88/oomstore/internal/database/metadata/postgres"

	"github.com/ethhte88/oomstore/internal/database/offline"
	offlinePG "github.com/ethhte88/oomstore/internal/database/offline/postgres"

	"github.com/ethhte88/oomstore/internal/database/online"
	onlinePG "github.com/ethhte88/oomstore/internal/database/online/postgres"
	onlineRedis "github.com/ethhte88/oomstore/internal/database/online/redis"
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
		return metadataPG.Open(context.Background(), opt.Postgres)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func CreateMetadataDatabase(ctx context.Context, opt types.MetadataStoreConfig) error {
	switch opt.Backend {
	case types.POSTGRES:
		return metadataPG.CreateDatabase(ctx, *opt.Postgres)
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
