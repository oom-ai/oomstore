package types

type BackendType string

const (
	POSTGRES BackendType = "postgres"
	REDIS    BackendType = "redis"
)

type OomStoreOpt struct {
	MetaStoreOpt    MetaStoreOpt    `yaml:"meta-store"`
	OfflineStoreOpt OfflineStoreOpt `yaml:"offline-store"`
	OnlineStoreOpt  OnlineStoreOpt  `yaml:"online-store"`
}

type OnlineStoreOpt struct {
	Backend       BackendType    `yaml:"backend"`
	PostgresDbOpt *PostgresDbOpt `yaml:"postgres"`
	RedisDbOpt    *RedisDbOpt    `yaml:"redis"`
}

type OfflineStoreOpt struct {
	Backend       BackendType    `yaml:"backend"`
	PostgresDbOpt *PostgresDbOpt `yaml:"postgres"`
}

type MetaStoreOpt struct {
	Backend       BackendType    `yaml:"backend"`
	PostgresDbOpt *PostgresDbOpt `yaml:"postgres"`
}

type RedisDbOpt struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Pass     string `yaml:"password"`
	Database int    `yaml:"database"`
}

type PostgresDbOpt struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Pass     string `yaml:"password"`
	Database string `yaml:"database"`
}
