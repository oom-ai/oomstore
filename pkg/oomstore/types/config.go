package types

import (
	"time"

	"github.com/oom-ai/oomstore/pkg/errdefs"
)

type BackendType string

const (
	BackendPostgres  BackendType = "postgres"
	BackendRedis     BackendType = "redis"
	BackendMySQL     BackendType = "mysql"
	BackendSQLite    BackendType = "sqlite"
	BackendSnowflake BackendType = "snowflake"
	BackendDynamoDB  BackendType = "dynamodb"
	BackendCassandra BackendType = "cassandra"
	BackendBigQuery  BackendType = "bigquery"
	BackendRedshift  BackendType = "redshift"
	BackendTiDB      BackendType = "tidb"
	BackendTiKV      BackendType = "tikv"
)

type OomStoreConfig struct {
	MetadataStore MetadataStoreConfig `yaml:"metadata-store"`
	OfflineStore  OfflineStoreConfig  `yaml:"offline-store"`
	OnlineStore   OnlineStoreConfig   `yaml:"online-store"`

	PushProcessor *PushProcessorConfig `yaml:"push-processor"`
}

type PushProcessorConfig struct {
	BufferSize    int           `yaml:"buffer-size"`
	FlushInterval time.Duration `yaml:"flush-interval"`
}

type OnlineStoreConfig struct {
	Backend   BackendType   `yaml:"-"`
	Postgres  *PostgresOpt  `yaml:"postgres"`
	Redis     *RedisOpt     `yaml:"redis"`
	MySQL     *MySQLOpt     `yaml:"mysql"`
	DynamoDB  *DynamoDBOpt  `yaml:"dynamodb"`
	Cassandra *CassandraOpt `yaml:"cassandra"`
	SQLite    *SQLiteOpt    `yaml:"sqlite"`
	TiDB      *MySQLOpt     `yaml:"tidb"`
	TiKV      *TiKVOpt      `yaml:"tikv"`
}

type OfflineStoreConfig struct {
	Backend   BackendType   `yaml:"-"`
	Postgres  *PostgresOpt  `yaml:"postgres"`
	MySQL     *MySQLOpt     `yaml:"mysql"`
	Snowflake *SnowflakeOpt `yaml:"snowflake"`
	BigQuery  *BigQueryOpt  `yaml:"bigquery"`
	Redshift  *RedshiftOpt  `yaml:"redshift"`
	SQLite    *SQLiteOpt    `yaml:"sqlite"`
	TiDB      *MySQLOpt     `yaml:"tidb"`
}

type MetadataStoreConfig struct {
	Backend  BackendType  `yaml:"-"`
	Postgres *PostgresOpt `yaml:"postgres"`
	MySQL    *MySQLOpt    `yaml:"mysql"`
	SQLite   *SQLiteOpt   `yaml:"sqlite"`
	TiDB     *MySQLOpt    `yaml:"tidb"`
}

type RedisOpt struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

type PostgresOpt struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type RedshiftOpt = PostgresOpt

type SQLiteOpt struct {
	DBFile string `yaml:"db-file"`
}

type MySQLOpt struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type SnowflakeOpt struct {
	Account  string `yaml:"account"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type BigQueryOpt struct {
	ProjectID   string `yaml:"project_id"`
	DatasetID   string `yaml:"dataset_id"`
	Credentials string `yaml:"credentials"`
}

type DynamoDBOpt struct {
	Region          string `yaml:"region"`
	EndpointURL     string `yaml:"endpoint-url"`
	AccessKeyID     string `yaml:"access-key-id"`
	SecretAccessKey string `yaml:"secret-access-key"`
	SessionToken    string `yaml:"session-token"`
	Source          string `yaml:"source"`
}

type CassandraOpt struct {
	Hosts    []string      `yaml:"hosts"`
	User     string        `yaml:"user"`
	Password string        `yaml:"password"`
	KeySpace string        `yaml:"keyspace"`
	Timeout  time.Duration `yaml:"timeout"`
}

type TiKVOpt struct {
	PdAddrs []string `yaml:"pd-addrs"`
}

func (cfg *OomStoreConfig) Validate() error {
	if err := cfg.MetadataStore.Validate(); err != nil {
		return err
	}
	if err := cfg.OnlineStore.Validate(); err != nil {
		return err
	}
	if err := cfg.OfflineStore.Validate(); err != nil {
		return err
	}
	return nil
}

func (cfg *MetadataStoreConfig) Validate() error {
	n := 0
	if cfg.Postgres != nil {
		cfg.Backend = BackendPostgres
		n++
	}
	if cfg.MySQL != nil {
		cfg.Backend = BackendMySQL
		n++
	}
	if cfg.SQLite != nil {
		cfg.Backend = BackendSQLite
		n++
	}
	if cfg.TiDB != nil {
		cfg.Backend = BackendTiDB
		n++
	}
	if n != 1 {
		return errdefs.Errorf("require exactly one metadata store backend")
	}
	return nil
}

func (cfg *OnlineStoreConfig) Validate() error {
	n := 0
	if cfg.Postgres != nil {
		cfg.Backend = BackendPostgres
		n++
	}
	if cfg.MySQL != nil {
		cfg.Backend = BackendMySQL
		n++
	}
	if cfg.Redis != nil {
		cfg.Backend = BackendRedis
		n++
	}
	if cfg.DynamoDB != nil {
		cfg.Backend = BackendDynamoDB
		n++
	}
	if cfg.Cassandra != nil {
		cfg.Backend = BackendCassandra
		n++
	}
	if cfg.SQLite != nil {
		cfg.Backend = BackendSQLite
		n++
	}
	if cfg.TiDB != nil {
		cfg.Backend = BackendTiDB
		n++
	}
	if cfg.TiKV != nil {
		cfg.Backend = BackendTiKV
		n++
	}
	if n != 1 {
		return errdefs.Errorf("require exactly one online store backend")
	}
	return nil
}

func (cfg *OfflineStoreConfig) Validate() error {
	n := 0
	if cfg.Postgres != nil {
		cfg.Backend = BackendPostgres
		n++
	}
	if cfg.MySQL != nil {
		cfg.Backend = BackendMySQL
		n++
	}
	if cfg.Snowflake != nil {
		cfg.Backend = BackendSnowflake
		n++
	}
	if cfg.BigQuery != nil {
		cfg.Backend = BackendBigQuery
		n++
	}
	if cfg.Redshift != nil {
		cfg.Backend = BackendRedshift
		n++
	}
	if cfg.SQLite != nil {
		cfg.Backend = BackendSQLite
		n++
	}
	if cfg.TiDB != nil {
		cfg.Backend = BackendTiDB
		n++
	}
	if n != 1 {
		return errdefs.Errorf("require exactly one offline store backend")
	}
	return nil
}
