package types

import (
	"fmt"
	"time"
)

type BackendType string

const (
	POSTGRES  BackendType = "postgres"
	REDIS     BackendType = "redis"
	MYSQL     BackendType = "mysql"
	SQLite    BackendType = "sqlite"
	SNOWFLAKE BackendType = "snowflake"
	DYNAMODB  BackendType = "dynamodb"
	CASSANDRA BackendType = "cassandra"
	BIGQUERY  BackendType = "bigquery"
	REDSHIFT  BackendType = "redshift"
)

type OomStoreConfig struct {
	MetadataStore MetadataStoreConfig `yaml:"metadata-store"`
	OfflineStore  OfflineStoreConfig  `yaml:"offline-store"`
	OnlineStore   OnlineStoreConfig   `yaml:"online-store"`
}

type OnlineStoreConfig struct {
	Backend   BackendType   `yaml:"-"`
	Postgres  *PostgresOpt  `yaml:"postgres"`
	Redis     *RedisOpt     `yaml:"redis"`
	MySQL     *MySQLOpt     `yaml:"mysql"`
	DynamoDB  *DynamoDBOpt  `yaml:"dynamodb"`
	Cassandra *CassandraOpt `yaml:"cassandra"`
}

type OfflineStoreConfig struct {
	Backend   BackendType   `yaml:"-"`
	Postgres  *PostgresOpt  `yaml:"postgres"`
	MySQL     *MySQLOpt     `yaml:"mysql"`
	Snowflake *SnowflakeOpt `yaml:"snowflake"`
	BigQuery  *BigQueryOpt  `yaml:"bigquery"`
	Redshift  *RedshiftOpt  `yaml:"redshift"`
}

type MetadataStoreConfig struct {
	Backend  BackendType  `yaml:"-"`
	Postgres *PostgresOpt `yaml:"postgres"`
	MySQL    *MySQLOpt    `yaml:"mysql"`
	SQLite   *SQLiteOpt   `yaml:"sqlite"`
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
		cfg.Backend = POSTGRES
		n++
	}
	if cfg.MySQL != nil {
		cfg.Backend = MYSQL
		n++
	}
	if n != 1 {
		return fmt.Errorf("require exactly one metadata store backend")
	}
	return nil
}

func (cfg *OnlineStoreConfig) Validate() error {
	n := 0
	if cfg.Postgres != nil {
		cfg.Backend = POSTGRES
		n++
	}
	if cfg.MySQL != nil {
		cfg.Backend = MYSQL
		n++
	}
	if cfg.Redis != nil {
		cfg.Backend = REDIS
		n++
	}
	if cfg.DynamoDB != nil {
		cfg.Backend = DYNAMODB
		n++
	}
	if n != 1 {
		return fmt.Errorf("require exactly one online store backend")
	}
	return nil
}

func (cfg *OfflineStoreConfig) Validate() error {
	n := 0
	if cfg.Postgres != nil {
		cfg.Backend = POSTGRES
		n++
	}
	if cfg.MySQL != nil {
		cfg.Backend = MYSQL
		n++
	}
	if cfg.Snowflake != nil {
		cfg.Backend = SNOWFLAKE
		n++
	}
	if n != 1 {
		return fmt.Errorf("require exactly one offline store backend")
	}
	return nil
}
