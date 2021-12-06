package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/ethhte88/oomstore/internal/database/metadata/sqlutil"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ethhte88/oomstore/internal/database/metadata"
	"github.com/ethhte88/oomstore/internal/database/metadata/informer"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

var _ metadata.Store = &DB{}
var _ metadata.DBStore = &Tx{}

type DB struct {
	*sqlx.DB
	*informer.Informer
}

type Tx struct {
	*sqlx.Tx
}

func Open(ctx context.Context, option *types.PostgresOpt) (*DB, error) {
	db, err := OpenDB(option.Host, option.Port, option.User, option.Password, option.Database)
	if err != nil {
		return nil, err
	}

	// TODO: make the interval configurable
	informer, err := informer.New(time.Second, func() (*informer.Cache, error) {
		return sqlutil.ListMetaData(ctx, db)
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return &DB{
		DB:       db,
		Informer: informer,
	}, nil
}

func OpenDB(host, port, user, password, database string) (*sqlx.DB, error) {
	return sqlx.Open(
		"postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			password,
			host,
			port,
			database),
	)
}

func (db *DB) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

func (db *DB) Close() error {
	if err := db.Informer.Close(); err != nil {
		return err
	}
	if err := db.DB.Close(); err != nil {
		return err
	}
	return nil
}

func (db *DB) WithTransaction(ctx context.Context, fn func(context.Context, metadata.DBStore) error) (err error) {
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

func (tx *Tx) WithTransaction(ctx context.Context, fn func(context.Context, metadata.DBStore) error) (err error) {
	return fn(ctx, tx)
}
