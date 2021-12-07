package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

var _ offline.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func (db *DB) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

func Open(option *types.PostgresOpt) (*DB, error) {
	db, err := dbutil.OpenPostgresDB(option.Host, option.Port, option.User, option.Password, option.Database)
	return &DB{db}, err
}
