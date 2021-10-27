package types

type BackendType string

const (
	POSTGRES BackendType = "postgres"
	REDIS    BackendType = "redis"
)

type OomStoreOptV2 struct {
	MetaStoreOpt    MetaStoreOpt
	OfflineStoreOpt OfflineStoreOpt
	OnlineStoreOpt  OnlineStoreOpt
}

type OnlineStoreOpt struct {
	Backend       BackendType
	PostgresDbOpt *PostgresDbOpt
	RedisDbOpt    *RedisDbOpt
}

type OfflineStoreOpt struct {
	Backend       BackendType
	PostgresDbOpt *PostgresDbOpt
}

type MetaStoreOpt struct {
	Backend       BackendType
	PostgresDbOpt *PostgresDbOpt
}

type RedisDbOpt struct {
	Host     string
	Port     string
	Pass     string
	Database int
}

type PostgresDbOpt struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
}
