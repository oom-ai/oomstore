package mysql

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/internal/database/offline/sqlutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
)

const MySQLBatchSize = 20

var _ offline.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func (db *DB) Ping(ctx context.Context) error {
	return db.PingContext(ctx)
}

func Open(option *types.MySQLOpt) (*DB, error) {
	db, err := dbutil.OpenMysqlDB(option.Host, option.Port, option.User, option.Password, option.Database)
	return &DB{db}, err
}

func (db *DB) Import(ctx context.Context, opt offline.ImportOpt) (int64, error) {
	return sqlutil.Import(ctx, db.DB, opt, dbutil.LoadDataFromCSVReader(types.MYSQL, MySQLBatchSize), types.MYSQL)
}

func (db *DB) Export(ctx context.Context, opt offline.ExportOpt) (<-chan types.ExportRecord, <-chan error) {
	return sqlutil.Export(ctx, db.DB, opt, types.MYSQL)
}

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	return sqlutil.Join(ctx, db.DB, opt, types.MYSQL)
}

func (db *DB) TypeTag(dbType string) (string, error) {
	return TypeTag(dbType)
}
