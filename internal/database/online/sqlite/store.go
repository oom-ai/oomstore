package sqlite

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const BackendType = types.BackendSQLite

var _ online.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func Open(opt *types.SQLiteOpt) (*DB, error) {
	db, err := dbutil.OpenSQLite(opt.DBFile)
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

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	panic("Implement me!")
}

func (db *DB) PrepareStreamTable(ctx context.Context, opt online.PrepareStreamTableOpt) error {
	tableName := sqlutil.OnlineStreamTableName(opt.GroupID)

	if opt.Feature == nil {
		schema, err := sqlutil.CreateStreamTableSchema(ctx, tableName, opt.Entity, types.BackendMySQL)
		if err != nil {
			return err
		}

		_, err = db.ExecContext(ctx, schema)
		return err
	}

	dbValueType, err := dbutil.DBValueType(types.BackendSQLite, opt.Feature.ValueType)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, opt.Feature.Name, dbValueType)
	_, err = db.ExecContext(ctx, sql)
	return err
}
