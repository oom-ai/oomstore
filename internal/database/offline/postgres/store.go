package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var PostgresTypeMap = map[string]types.ValueType{
	"bigint":    types.INT64,
	"int8":      types.INT64,
	"bigserial": types.INT64,
	"serial8":   types.INT64,

	"boolean": types.BOOL,
	"bool":    types.BOOL,

	"bytea":       types.BYTES,
	"jsonb":       types.BYTES,
	"uuid":        types.BYTES,
	"bit":         types.BYTES,
	"bit varying": types.BYTES,
	"character":   types.BYTES,
	"char":        types.BYTES,
	"json":        types.BYTES,
	"money":       types.BYTES,
	"numeric":     types.BYTES,

	"character varying": types.STRING,
	"text":              types.STRING,
	"varchar":           types.STRING,

	"double precision": types.FLOAT64,
	"float8":           types.FLOAT64,

	"integer": types.INT64,
	"int":     types.INT64,
	"int4":    types.INT64,
	"serial":  types.INT64,
	"serial4": types.INT64,

	"real":   types.FLOAT64,
	"float4": types.FLOAT64,

	"smallint":    types.INT64,
	"int2":        types.INT64,
	"smallserial": types.INT64,
	"serial2":     types.INT64,

	"date":                        types.TIME,
	"time":                        types.TIME,
	"time without time zone":      types.TIME,
	"time with time zone":         types.TIME,
	"timetz":                      types.TIME,
	"timestamp":                   types.TIME,
	"timestamp without time zone": types.TIME,
	"timestamp with time zone":    types.TIME,
	"timestamptz":                 types.TIME,
}

var _ offline.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func (db *DB) Ping(ctx context.Context) error {
	return db.PingContext(ctx)
}

func Open(option *types.PostgresOpt) (*DB, error) {
	db, err := dbutil.OpenPostgresDB(option.Host, option.Port, option.User, option.Password, option.Database)
	return &DB{db}, err
}

func (db *DB) Import(ctx context.Context, opt offline.ImportOpt) (int64, error) {
	return sqlutil.Import(ctx, db.DB, opt, loadDataFromSource, types.POSTGRES)
}

func (db *DB) Export(ctx context.Context, opt offline.ExportOpt) (<-chan types.ExportRecord, <-chan error) {
	return sqlutil.Export(ctx, db.DB, opt, types.POSTGRES)
}

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	return sqlutil.Join(ctx, db.DB, opt, types.POSTGRES)
}

func (db *DB) TypeTag(dbType string) (types.ValueType, error) {
	return sqlutil.GetValueType(PostgresTypeMap, dbType)
}

func (db *DB) TableSchema(ctx context.Context, tableName string) (*types.DataTableSchema, error) {
	rows, err := db.QueryxContext(ctx, "select column_name, data_type from information_schema.columns where table_name = $1", tableName)
	if err != nil {
		return nil, err
	}
	return sqlutil.SqlxTableSchema(ctx, db, rows)
}
