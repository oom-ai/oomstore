package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/internal/database/online/sqlutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

const BackendType = types.POSTGRES

var _ online.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func Open(opt *types.PostgresOpt) (*DB, error) {
	db, err := dbutil.OpenPostgresDB(opt.Host, opt.Port, opt.User, opt.Password, opt.Database)
	return &DB{db}, err
}

func (db *DB) Ping(ctx context.Context) error {
	return db.PingContext(ctx)
}

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	return sqlutil.Get(ctx, db.DB, opt, BackendType)
}

func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	return sqlutil.MultiGet(ctx, db.DB, opt, BackendType)
}

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	return sqlutil.Import(ctx, db.DB, opt, BackendType)
}

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	return sqlutil.Purge(ctx, db.DB, revisionID, BackendType)
}
