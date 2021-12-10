package types

import "time"

type BackendType string

const (
	POSTGRES  BackendType = "postgres"
	REDIS     BackendType = "redis"
	MYSQL     BackendType = "mysql"
	SNOWFLAKE BackendType = "snowflake"
	DYNAMODB  BackendType = "dynamodb"
	CASSANDRA BackendType = "cassandra"
)

type OomStoreConfig struct {
	MetadataStore MetadataStoreConfig `yaml:"metadata-store"`
	OfflineStore  OfflineStoreConfig  `yaml:"offline-store"`
	OnlineStore   OnlineStoreConfig   `yaml:"online-store"`
}

type OnlineStoreConfig struct {
	Backend  BackendType  `yaml:"backend"`
	Postgres *PostgresOpt `yaml:"postgres"`
	Redis    *RedisOpt    `yaml:"redis"`
	MySQL    *MySQLOpt    `yaml:"mysql"`
	DynamoDB *DynamoDBOpt `yaml:"dynamodb"`
}

type OfflineStoreConfig struct {
	Backend   BackendType   `yaml:"backend"`
	Postgres  *PostgresOpt  `yaml:"postgres"`
	MySQL     *MySQLOpt     `yaml:"mysql"`
	Snowflake *SnowflakeOpt `yaml:"snowflake"`
}

type MetadataStoreConfig struct {
	Backend  BackendType  `yaml:"backend"`
	Postgres *PostgresOpt `yaml:"postgres"`
	MySQL    *MySQLOpt    `yaml:"mysql"`
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
