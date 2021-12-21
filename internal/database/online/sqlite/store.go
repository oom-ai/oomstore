package sqlite

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const BackendType = types.SQLite

var _ online.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func Open(opt *types.SQLiteOpt) (*DB, error) {
	db, err := dbutil.OpenSQLite(opt.DBFile)
	return &DB{db}, err
}

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	//TODO implement me
	panic("implement me")
}

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	//TODO implement me
	panic("implement me")
}

func (db *DB) Ping(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
