package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type DB struct {
	*sqlx.DB
}

func Open(option *types.PostgresDbOpt) (*DB, error) {
	return OpenWith(option.Host, option.Port, option.User, option.Pass, option.Database)
}

func OpenWith(host, port, user, pass, database string) (*DB, error) {
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			pass,
			host,
			port,
			database),
	)
	return &DB{db}, err
}
