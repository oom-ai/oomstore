package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata/informer"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

//var _ metadata.Store = &DB{}
//var _ metadata.DBStore = &Tx{}

type DB struct {
	*sqlx.DB
	*informer.Informer
}

type Tx struct {
	*sqlx.Tx
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

func Open(ctx context.Context, option *types.MySQLOpt) (*DB, error) {
	db, err := OpenDB(option.Host, option.Port, option.User, option.Password, option.Database)
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
	return &DB{
		DB:       db,
		Informer: informer,
	}, nil
}

func OpenDB(host, port, user, password, database string) (*sqlx.DB, error) {
	return sqlx.Open(
		"mysql",
		fmt.Sprintf("%s:%s@(%s:%s)/%s", user, password, host, port, database))
}

func CreateDatabase(ctx context.Context, opt types.MySQLOpt) (err error) {
	defaultDB, err := OpenDB(opt.Host, opt.Port, opt.User, opt.Password, "")
	if err != nil {
		return
	}
	defer defaultDB.Close()

	if _, err = defaultDB.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", opt.Database)); err != nil {
		return
	}

	db, err := OpenDB(opt.Host, opt.Port, opt.User, opt.Password, opt.Database)
	if err != nil {
		return
	}
	defer db.Close()
	return createMetaSchemas(ctx, db)
}

func createMetaSchemas(ctx context.Context, db *sqlx.DB) (err error) {
	// Use transaction to guarantee the following operations be executed
	// on the same connection: http://go-database-sql.org/modifying.html
	return dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create meta tables
		for _, schema := range META_TABLE_SCHEMAS {
			if _, err = tx.ExecContext(ctx, schema); err != nil {
				return fmt.Errorf("failed to create table, err=%+v", err)
			}
		}

		// create foreign keys
		for _, stmt := range META_TABLE_FOREIGN_KEYS {
			if _, err = tx.ExecContext(ctx, stmt); err != nil {
				return fmt.Errorf("failed to add foreign key, err=%+v", err)
			}
		}

		// create meta views
		for _, schema := range META_VIEW_SCHEMAS {
			if _, err = tx.ExecContext(ctx, schema); err != nil {
				return err
			}
		}

		return nil
	})
}

func list(ctx context.Context, db *sqlx.DB) (*informer.Cache, error) {
	var cache *informer.Cache
	err := dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		entities := types.EntityList{}
		if err := tx.SelectContext(ctx, &entities, `SELECT * FROM entity`); err != nil {
			return err
		}

		features := types.FeatureList{}
		if err := tx.SelectContext(ctx, &features, `SELECT * FROM feature`); err != nil {
			return err
		}

		groups := types.GroupList{}
		if err := tx.SelectContext(ctx, &groups, `SELECT * FROM feature_group`); err != nil {
			return err
		}

		cache = informer.NewCache(entities, features, groups)
		return nil
	})
	return cache, err
}
