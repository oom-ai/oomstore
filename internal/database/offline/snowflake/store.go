package snowflake

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/snowflakedb/gosnowflake"
)

const SnowflakeBatchSize = 100

// TODO: add type NUMBER, DECIMAL, NUMERIC
var SnowflakeTypeMap = map[string]types.ValueType{
	"boolean": types.BOOL,

	"binary":    types.BYTES,
	"varbinary": types.BYTES,

	"integer":  types.INT64,
	"int":      types.INT64,
	"smallint": types.INT64,
	"bigint":   types.INT64,
	"tinyint":  types.INT64,
	"byteint":  types.INT64,

	"double":           types.FLOAT64,
	"double precision": types.FLOAT64,
	"real":             types.FLOAT64,
	"float":            types.FLOAT64,
	"float4":           types.FLOAT64,
	"float8":           types.FLOAT64,

	"string":    types.STRING,
	"text":      types.STRING,
	"varchar":   types.STRING,
	"char":      types.STRING,
	"character": types.STRING,

	"date":      types.TIME,
	"time":      types.TIME,
	"datetime":  types.TIME,
	"timestamp": types.TIME,
}

var _ offline.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func (db *DB) Ping(ctx context.Context) error {
	return db.PingContext(ctx)
}

func Open(opt *types.SnowflakeOpt) (*DB, error) {
	dsn, err := gosnowflake.DSN(&gosnowflake.Config{
		Account:  opt.Account,
		User:     opt.User,
		Password: opt.Password,
		Database: opt.Database,
	})
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("snowflake", dsn)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, err
}

func (db *DB) Import(ctx context.Context, opt offline.ImportOpt) (int64, error) {
	return sqlutil.Import(ctx, db.DB, opt, dbutil.LoadDataFromSource(types.SNOWFLAKE, SnowflakeBatchSize), types.SNOWFLAKE)
}

func (db *DB) Export(ctx context.Context, opt offline.ExportOpt) (<-chan types.ExportRecord, <-chan error) {
	return sqlutil.Export(ctx, db.DB, opt, types.SNOWFLAKE)
}

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	return sqlutil.Join(ctx, db.DB, opt, types.SNOWFLAKE)
}

func (db *DB) TypeTag(dbType string) (types.ValueType, error) {
	return sqlutil.GetValueType(SnowflakeTypeMap, dbType)
}

func (db *DB) TableSchema(ctx context.Context, tableName string) (*types.DataTableSchema, error) {
	return nil, fmt.Errorf("not implemented")
}
