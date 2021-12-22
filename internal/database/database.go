package database

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	metadataMySQL "github.com/oom-ai/oomstore/internal/database/metadata/mysql"
	metadataPG "github.com/oom-ai/oomstore/internal/database/metadata/postgres"
	metadataSQLite "github.com/oom-ai/oomstore/internal/database/metadata/sqlite"

	"github.com/oom-ai/oomstore/internal/database/offline"
	offlineBigQuery "github.com/oom-ai/oomstore/internal/database/offline/bigquery"
	offlineMySQL "github.com/oom-ai/oomstore/internal/database/offline/mysql"
	offlinePG "github.com/oom-ai/oomstore/internal/database/offline/postgres"
	offlineRedshift "github.com/oom-ai/oomstore/internal/database/offline/redshift"
	offlineSnowflake "github.com/oom-ai/oomstore/internal/database/offline/snowflake"

	"github.com/oom-ai/oomstore/internal/database/online"
	onlineCassandra "github.com/oom-ai/oomstore/internal/database/online/cassandra"
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
	case types.CASSANDRA:
		return onlineCassandra.Open(opt.Cassandra)
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
	case types.SQLite:
		return metadataSQLite.Open(context.Background(), opt.SQLite)
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
	case types.SQLite:
		return metadataSQLite.CreateDatabase(ctx, *opt.SQLite)
	default:
		return fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func OpenOfflineStore(ctx context.Context, opt types.OfflineStoreConfig) (offline.Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return offlinePG.Open(opt.Postgres)
	case types.MYSQL:
		return offlineMySQL.Open(opt.MySQL)
	case types.SNOWFLAKE:
		return offlineSnowflake.Open(opt.Snowflake)
	case types.BIGQUERY:
		return offlineBigQuery.Open(ctx, opt.BigQuery)
	case types.REDSHIFT:
		return offlineRedshift.Open(opt.Redshift)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}
