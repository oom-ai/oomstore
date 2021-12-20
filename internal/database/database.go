package database

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	metadataMySQL "github.com/oom-ai/oomstore/internal/database/metadata/mysql"
	metadataPG "github.com/oom-ai/oomstore/internal/database/metadata/postgres"

	"github.com/oom-ai/oomstore/internal/database/offline"
	offlineMySQL "github.com/oom-ai/oomstore/internal/database/offline/mysql"
	offlinePG "github.com/oom-ai/oomstore/internal/database/offline/postgres"
	offlineSnowflake "github.com/oom-ai/oomstore/internal/database/offline/snowflake"

	"github.com/oom-ai/oomstore/internal/database/online"
	onlineDynamoDB "github.com/oom-ai/oomstore/internal/database/online/dynamodb"
	onlineMySQL "github.com/oom-ai/oomstore/internal/database/online/mysql"
	onlinePG "github.com/oom-ai/oomstore/internal/database/online/postgres"
	onlineRedis "github.com/oom-ai/oomstore/internal/database/online/redis"
)

func OpenOnlineStore(opt types.OnlineStoreConfig) (online.Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return onlinePG.Open(opt.Postgres)
	case types.REDIS:
		return onlineRedis.Open(opt.Redis), nil
	case types.MYSQL:
		return onlineMySQL.Open(opt.MySQL)
	case types.DYNAMODB:
		return onlineDynamoDB.Open(opt.DynamoDB)
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
