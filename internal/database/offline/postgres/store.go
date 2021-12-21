package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const PostgresBatchSize = 20

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
	return sqlutil.Import(ctx, db.DB, opt, dbutil.LoadDataFromSource(types.POSTGRES, PostgresBatchSize), types.POSTGRES)
}

func (db *DB) Export(ctx context.Context, opt offline.ExportOpt) (<-chan types.ExportRecord, <-chan error) {
	return sqlutil.Export(ctx, db.DB, opt, types.POSTGRES)
}

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	return sqlutil.Join(ctx, db.DB, opt, types.POSTGRES)
}

func (db *DB) TypeTag(dbType string) (string, error) {
	return TypeTag(dbType)
}

func (db *DB) TableSchema(ctx context.Context, tableName string) (*types.DataTableSchema, error) {
	rows, err := db.QueryxContext(ctx, "select column_name, data_type from information_schema.columns where table_name = $1", tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schema types.DataTableSchema
	for rows.Next() {
		var fieldName, dbValueType string
		if err := rows.Scan(&fieldName, &dbValueType); err != nil {
			return nil, err
		}
		valueType, err := db.TypeTag(dbValueType)
		if err != nil {
			return nil, err
		}
		schema.Fields = append(schema.Fields, types.DataTableFieldSchema{
			Name:      fieldName,
			ValueType: valueType,
		})
	}
	if len(schema.Fields) == 0 {
		return nil, fmt.Errorf("table not found")
	}
	return &schema, nil
}
