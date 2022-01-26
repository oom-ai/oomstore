package postgres

import (
	"context"
	"fmt"

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

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	tableName := sqlutil.OnlineStreamTableName(opt.GroupID)

	cond := sqlutil.BuildPushCondition(opt, Backend)

	query := fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s) ON CONFLICT ("%s") DO UPDATE SET %s`,
		tableName,
		cond.Inserts,
		cond.InsertPlaceholders,
		opt.EntityName,
		cond.UpdatePlaceholders,
	)

	_, err := db.ExecContext(ctx, db.Rebind(query), append(cond.InsertValues, cond.UpdateValues...)...)
	return errdefs.WithStack(err)
}

func (db *DB) CreateTable(ctx context.Context, opt online.CreateTableOpt) error {
	dbOpt := dbutil.DBOpt{
		Backend: Backend,
		SqlxDB:  db.DB,
	}
	return sqlutil.CreateTable(ctx, dbOpt, opt)
}
