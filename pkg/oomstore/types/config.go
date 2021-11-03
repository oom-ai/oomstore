package types

type BackendType string

const (
	POSTGRES BackendType = "postgres"
	REDIS    BackendType = "redis"
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
}

type OfflineStoreConfig struct {
	Backend  BackendType  `yaml:"backend"`
	Postgres *PostgresOpt `yaml:"postgres"`
}

type MetadataStoreConfig struct {
	Backend  BackendType  `yaml:"backend"`
	Postgres *PostgresOpt `yaml:"postgres"`
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
