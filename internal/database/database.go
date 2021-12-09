package database

import (
	"context"
	"fmt"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"

	"github.com/ethhte88/oomstore/internal/database/metadata"
	metadataMySQL "github.com/ethhte88/oomstore/internal/database/metadata/mysql"
	metadataPG "github.com/ethhte88/oomstore/internal/database/metadata/postgres"

	"github.com/ethhte88/oomstore/internal/database/offline"
	offlineMySQL "github.com/ethhte88/oomstore/internal/database/offline/mysql"
	offlinePG "github.com/ethhte88/oomstore/internal/database/offline/postgres"
	offlineSnowflake "github.com/ethhte88/oomstore/internal/database/offline/snowflake"

	"github.com/ethhte88/oomstore/internal/database/online"
	onlineMySQL "github.com/ethhte88/oomstore/internal/database/online/mysql"
	onlinePG "github.com/ethhte88/oomstore/internal/database/online/postgres"
	onlineRedis "github.com/ethhte88/oomstore/internal/database/online/redis"
)

func OpenOnlineStore(opt types.OnlineStoreConfig) (online.Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return onlinePG.Open(opt.Postgres)
	case types.REDIS:
		return onlineRedis.Open(opt.Redis), nil
	case types.MYSQL:
		return onlineMySQL.Open(opt.MySQL)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func OpenMetadataStore(opt types.MetadataStoreConfig) (metadata.Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return metadataPG.Open(context.Background(), opt.Postgres)
	case types.MYSQL:
		return metadataMySQL.Open(context.Background(), opt.MySQL)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func CreateMetadataDatabase(ctx context.Context, opt types.MetadataStoreConfig) error {
	switch opt.Backend {
	case types.POSTGRES:
		return metadataPG.CreateDatabase(ctx, *opt.Postgres)
	case types.MYSQL:
		return metadataMySQL.CreateDatabase(ctx, *opt.MySQL)
	default:
		return fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func OpenOfflineStore(opt types.OfflineStoreConfig) (offline.Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return offlinePG.Open(opt.Postgres)
	case types.MYSQL:
		return offlineMySQL.Open(opt.MySQL)
	case types.SNOWFLAKE:
		return offlineSnowflake.Open(opt.Snowflake)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}
