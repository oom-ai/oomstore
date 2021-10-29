package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var _ offline.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func Open(option *types.PostgresOpt) (*DB, error) {
	return OpenWith(option.Host, option.Port, option.User, option.Password, option.Database)
}

func OpenWith(host, port, user, password, database string) (*DB, error) {
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			password,
			host,
			port,
			database),
	)
	return &DB{db}, err
}
