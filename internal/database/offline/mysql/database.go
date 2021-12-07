package mysql

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
)

var _ offline.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func (db *DB) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

func Open(option *types.MySQLOpt) (*DB, error) {
	return OpenWith(option.Host, option.Port, option.User, option.Password, option.Database)
}

func OpenWith(host, port, user, password, database string) (*DB, error) {
	db, err := sqlx.Open("mysql", dbutil.GetMysqlDSN(host, port, user, password, database))
	return &DB{
		DB: db,
	}, err
}
