package mysql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	Backend        = types.BackendMySQL
	MySQLBatchSize = 20
)

var _ offline.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func (db *DB) Ping(ctx context.Context) error {
	err := db.PingContext(ctx)
	return errdefs.WithStack(err)
}

func Open(option *types.MySQLOpt) (*DB, error) {
	db, err := dbutil.OpenMysqlDB(option.Host, option.Port, option.User, option.Password, option.Database)
	return &DB{db}, err
}

func (db *DB) Import(ctx context.Context, opt offline.ImportOpt) (int64, error) {
	return sqlutil.Import(ctx, db.DB, opt, dbutil.LoadDataFromSource(Backend, MySQLBatchSize), Backend)
}

func (db *DB) Export(ctx context.Context, opt offline.ExportOpt) (*types.ExportResult, error) {
	return sqlutil.Export(ctx, db.DB, opt, Backend)
}

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	return sqlutil.Join(ctx, db.DB, opt, Backend)
}

func (db *DB) TableSchema(ctx context.Context, opt offline.TableSchemaOpt) (*types.DataTableSchema, error) {
	rows, err := db.QueryxContext(ctx, "SELECT column_name, data_type FROM information_schema.columns WHERE table_name = ?", opt.TableName)
	if err != nil {
		return nil, err
	}
	return sqlutil.SqlxTableSchema(ctx, db.DB, Backend, rows, opt)
}

func (db *DB) Snapshot(ctx context.Context, opt offline.SnapshotOpt) error {
	dbOpt := dbutil.DBOpt{Backend: Backend, SqlxDB: db.DB}
	return sqlutil.Snapshot(ctx, dbOpt, opt)
}

func (db *DB) CreateTable(ctx context.Context, opt offline.CreateTableOpt) error {
	dbOpt := dbutil.DBOpt{Backend: Backend, SqlxDB: db.DB}
	return sqlutil.CreateTable(ctx, dbOpt, opt)
}

func (db *DB) Push(ctx context.Context, opt offline.PushOpt) error {
	dbOpt := dbutil.DBOpt{Backend: Backend, SqlxDB: db.DB}
	if err := sqlutil.Push(ctx, dbOpt, opt); err != nil {
		return err
	}
	return nil
}
