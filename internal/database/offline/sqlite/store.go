package sqlite

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
	Backend         = types.BackendSQLite
	SQLiteBatchSize = 20
)

var _ offline.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func (db *DB) Ping(ctx context.Context) error {
	err := db.PingContext(ctx)
	return errdefs.WithStack(err)
}

func Open(option *types.SQLiteOpt) (*DB, error) {
	db, err := dbutil.OpenSQLite(option.DBFile)
	return &DB{db}, errdefs.WithStack(err)
}

func (db *DB) Import(ctx context.Context, opt offline.ImportOpt) (int64, error) {
	return sqlutil.Import(ctx, db.DB, opt, dbutil.LoadDataFromSource(Backend, SQLiteBatchSize), Backend)
}

func (db *DB) ExportOneGroup(ctx context.Context, opt offline.ExportOneGroupOpt) (<-chan types.ExportRecord, <-chan error) {
	return sqlutil.ExportOneGroup(ctx, db.DB, opt, Backend)
}

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	return sqlutil.Join(ctx, db.DB, opt, Backend)
}

func (db *DB) TableSchema(ctx context.Context, tableName string) (*types.DataTableSchema, error) {
	query := `
		SELECT
			p.name as column_name,
			p.type AS data_type
		FROM sqlite_master AS m
		LEFT OUTER JOIN pragma_table_info((m.name)) AS p
		WHERE m.type = 'table' AND m.name = ?
`
	rows, err := db.QueryxContext(ctx, query, tableName)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	return sqlutil.SqlxTableSchema(ctx, db, types.BackendSQLite, rows)
}

func (db *DB) Snapshot(ctx context.Context, opt offline.SnapshotOpt) error {
	dbOpt := dbutil.DBOpt{Backend: Backend, SqlxDB: db.DB}
	return sqlutil.Snapshot(ctx, dbOpt, opt)
}

func (db *DB) CreateTable(ctx context.Context, opt offline.CreateTableOpt) error {
	return sqlutil.CreateTable(ctx, db.DB, opt, Backend)
}

func (db *DB) Push(ctx context.Context, opt offline.PushOpt) error {
	dbOpt := dbutil.DBOpt{Backend: Backend, SqlxDB: db.DB}
	if err := sqlutil.Push(ctx, dbOpt, opt); err != nil {
		return err
	}
	return nil
}
