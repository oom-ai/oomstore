package snowflake

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
	"github.com/snowflakedb/gosnowflake"
)

var _ offline.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func (db *DB) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

func Open(option *types.SnowflakeOpt) (*DB, error) {
	return OpenWith(option.Account, option.User, option.Password, option.Database)
}

func OpenWith(account, user, password, database string) (*DB, error) {
	dsn, err := gosnowflake.DSN(&gosnowflake.Config{
		Account:  account,
		User:     user,
		Password: password,
		Database: database,
	})
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("snowflake", dsn)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, err
}
