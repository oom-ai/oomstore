package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var _ metadata.Store = &DB{}

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

type Tx struct {
	*sqlx.Tx
}

func openTx(ctx context.Context, db *sqlx.DB) (*Tx, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx}, nil
}
