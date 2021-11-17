package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/informer"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var _ metadata.Store = &DB{}
var _ metadata.WriteStore = &Tx{}

type DB struct {
	*sqlx.DB
	*informer.Informer
}

type Tx struct {
	*sqlx.Tx
}

func Open(ctx context.Context, option *types.PostgresOpt) (*DB, error) {
	db, err := OpenDB(ctx, option.Host, option.Port, option.User, option.Password, option.Database)
	if err != nil {
		return nil, err
	}

	// TODO: make the interval configurable
	informer, err := informer.New(time.Second, func() (*informer.Cache, error) {
		return list(ctx, db)
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return &DB{DB: db, Informer: informer}, nil
}

func OpenDB(ctx context.Context, host, port, user, password, database string) (*sqlx.DB, error) {
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

func (db *DB) Close() error {
	if err := db.Informer.Close(); err != nil {
		return err
	}
	if err := db.DB.Close(); err != nil {
		return err
	}
	return nil
}

func list(ctx context.Context, db *sqlx.DB) (*informer.Cache, error) {
	var cache *informer.Cache
	err := dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		entities := types.EntityList{}
		if err := tx.SelectContext(ctx, &entities, `SELECT * FROM "feature_entity"`); err != nil {
			return err
		}

		features := types.FeatureList{}
		if err := tx.SelectContext(ctx, &features, `SELECT * FROM "feature"`); err != nil {
			return err
		}

		groups := types.FeatureGroupList{}
		if err := tx.SelectContext(ctx, &groups, `SELECT * FROM "feature_group"`); err != nil {
			return err
		}

		revisions := types.RevisionList{}
		if err := tx.SelectContext(ctx, &revisions, `SELECT * FROM "feature_group_revision"`); err != nil {
			return err
		}
		cache = informer.NewCache(entities, features, groups, revisions)
		return nil
	})
	return cache, err
}

func (db *DB) WithTransaction(ctx context.Context, fn func(context.Context, metadata.WriteStore) error) (err error) {
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

func (tx *Tx) WithTransaction(ctx context.Context, fn func(context.Context, metadata.WriteStore) error) (err error) {
	return fn(ctx, tx)
}
