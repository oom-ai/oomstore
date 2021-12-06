package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

var _ online.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func (db *DB) Ping(ctx context.Context) error {
	return db.DB.Ping()
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

func Open(opt *types.PostgresOpt) (*DB, error) {
	return OpenWith(opt.Host, opt.Port, opt.User, opt.Password, opt.Database)
}
