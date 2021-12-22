package mysql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const MySQLBatchSize = 20

var MySQLTypeMap = map[string]types.ValueType{
	"boolean": types.BOOL,
	"bool":    types.BOOL,

	"binary":    types.BYTES,
	"varbinary": types.BYTES,

	"integer":   types.INT64,
	"int":       types.INT64,
	"smallint":  types.INT64,
	"bigint":    types.INT64,
	"tinyint":   types.INT64,
	"mediumint": types.INT64,

	"double": types.FLOAT64,
	"float":  types.FLOAT64,

	"text":    types.STRING,
	"varchar": types.STRING,
	"char":    types.STRING,

	"date":      types.TIME,
	"time":      types.TIME,
	"datetime":  types.TIME,
	"timestamp": types.TIME,
	"year":      types.TIME,
}

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
	return sqlutil.Import(ctx, db.DB, opt, dbutil.LoadDataFromSource(types.MYSQL, MySQLBatchSize), types.MYSQL)
}

func (db *DB) Export(ctx context.Context, opt offline.ExportOpt) (<-chan types.ExportRecord, <-chan error) {
	return sqlutil.Export(ctx, db.DB, opt, types.MYSQL)
}

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	return sqlutil.Join(ctx, db.DB, opt, types.MYSQL)
}

func (db *DB) ValueType(dbType string) (types.ValueType, error) {
	return sqlutil.GetValueType(MySQLTypeMap, dbType)
}

func (db *DB) TableSchema(ctx context.Context, tableName string) (*types.DataTableSchema, error) {
	return nil, fmt.Errorf("not implemented")
}
