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
	offlineSQLite "github.com/oom-ai/oomstore/internal/database/offline/sqlite"

	"github.com/oom-ai/oomstore/internal/database/online"
	onlineCassandra "github.com/oom-ai/oomstore/internal/database/online/cassandra"
	onlineDynamoDB "github.com/oom-ai/oomstore/internal/database/online/dynamodb"
	onlineMySQL "github.com/oom-ai/oomstore/internal/database/online/mysql"
	onlinePG "github.com/oom-ai/oomstore/internal/database/online/postgres"
	onlineRedis "github.com/oom-ai/oomstore/internal/database/online/redis"
	onlineSQLite "github.com/oom-ai/oomstore/internal/database/online/sqlite"
	onlineTiKV "github.com/oom-ai/oomstore/internal/database/online/tikv"
)

func OpenOnlineStore(opt types.OnlineStoreConfig) (online.Store, error) {
	switch opt.Backend {
	case types.BackendPostgres:
		return onlinePG.Open(opt.Postgres)
	case types.BackendRedis:
		return onlineRedis.Open(opt.Redis), nil
	case types.BackendMySQL:
		return onlineMySQL.Open(opt.MySQL)
	case types.BackendTiDB:
		return onlineMySQL.Open(opt.TiDB)
	case types.BackendDynamoDB:
		return onlineDynamoDB.Open(opt.DynamoDB)
	case types.BackendCassandra:
		return onlineCassandra.Open(opt.Cassandra)
	case types.BackendSQLite:
		return onlineSQLite.Open(opt.SQLite)
	case types.BackendTiKV:
		return onlineTiKV.Open(opt.TiKV)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func OpenMetadataStore(opt types.MetadataStoreConfig) (metadata.Store, error) {
	switch opt.Backend {
	case types.BackendPostgres:
		return metadataPG.Open(context.Background(), opt.Postgres)
	case types.BackendMySQL:
		return metadataMySQL.Open(context.Background(), opt.MySQL)
	case types.BackendTiDB:
		return metadataMySQL.Open(context.Background(), opt.TiDB)
	case types.BackendSQLite:
		return metadataSQLite.Open(context.Background(), opt.SQLite)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func CreateMetadataDatabase(ctx context.Context, opt types.MetadataStoreConfig) error {
	switch opt.Backend {
	case types.BackendPostgres:
		return metadataPG.CreateDatabase(ctx, *opt.Postgres)
	case types.BackendMySQL:
		return metadataMySQL.CreateDatabase(ctx, *opt.MySQL)
	case types.BackendTiDB:
		return metadataMySQL.CreateDatabase(ctx, *opt.TiDB)
	case types.BackendSQLite:
		return metadataSQLite.CreateDatabase(ctx, *opt.SQLite)
	default:
		return fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}

func OpenOfflineStore(ctx context.Context, opt types.OfflineStoreConfig) (offline.Store, error) {
	switch opt.Backend {
	case types.BackendPostgres:
		return offlinePG.Open(opt.Postgres)
	case types.BackendMySQL:
		return offlineMySQL.Open(opt.MySQL)
	case types.BackendTiDB:
		return offlineMySQL.Open(opt.TiDB)
	case types.BackendSnowflake:
		return offlineSnowflake.Open(opt.Snowflake)
	case types.BackendBigQuery:
		return offlineBigQuery.Open(ctx, opt.BigQuery)
	case types.BackendRedshift:
		return offlineRedshift.Open(opt.Redshift)
	case types.BackendSQLite:
		return offlineSQLite.Open(opt.SQLite)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}
