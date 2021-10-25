package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type DB struct {
	*sqlx.DB
}

func OpenWith(host string, port, user, pass, database string) (*DB, error) {
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

func Open(opt *types.PostgresDbOpt) (*DB, error) {
	return OpenWith(opt.Host, opt.Port, opt.User, opt.Pass, opt.Database)
}
