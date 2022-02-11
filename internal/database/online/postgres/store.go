package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	Backend = types.BackendPostgres
)

var _ online.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func Open(opt *types.PostgresOpt) (*DB, error) {
	db, err := dbutil.OpenPostgresDB(opt.Host, opt.Port, opt.User, opt.Password, opt.Database)
	return &DB{db}, err
}

func (db *DB) Ping(ctx context.Context) error {
	return errdefs.WithStack(db.PingContext(ctx))
}

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	return sqlutil.Get(ctx, db.DB, opt, Backend)
}

func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	return sqlutil.MultiGet(ctx, db.DB, opt, Backend)
}

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	return sqlutil.Import(ctx, db.DB, opt, Backend)
}

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	return sqlutil.Purge(ctx, db.DB, revisionID, Backend)
}

const PUSH_QUERY = `
INSERT INTO {{ .TableName }} ( {{ .Fields }} )
VALUES ( {{ .InsertPlaceholders }} )
ON CONFLICT ( {{ qt .EntityName }} )
DO UPDATE SET {{ .UpdatePlaceholders }}
`

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	params := sqlutil.BuildPushQueryParams(opt, Backend)
	query, err := sqlutil.BuildPushQuery(params, PUSH_QUERY)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, db.Rebind(query), append(params.InsertValues, params.UpdateValues...)...)
	return errdefs.WithStack(err)
}

func (db *DB) CreateTable(ctx context.Context, opt online.CreateTableOpt) error {
	dbOpt := dbutil.DBOpt{
		Backend: Backend,
		SqlxDB:  db.DB,
	}
	return sqlutil.CreateTable(ctx, dbOpt, opt)
}
