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
var _ metadata.Store = &Tx{}

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

func (tx *Tx) Close() error {
	return nil
}

func (db *DB) WithTransaction(ctx context.Context, fn metadata.TxFn) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}
	txStore := &Tx{Tx: tx}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	return fn(ctx, txStore)
}

func (tx *Tx) WithTransaction(ctx context.Context, fn metadata.TxFn) (err error) {
	return fn(ctx, tx)
}
