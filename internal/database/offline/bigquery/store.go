package bigquery

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

var _ offline.Store = &DB{}

type DB struct {
	*bigquery.Client
}

func Open(ctx context.Context, opt *types.BigQueryOpt) (*DB, error) {
	client, err := bigquery.NewClient(ctx, opt.ProjectID)
	if err != nil {
		return nil, err
	}
	return &DB{client}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	q := db.Client.Query("SELECT 1")
	_, err := q.Read(ctx)
	return err
}

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	panic("implement me")
}

func (db *DB) Export(ctx context.Context, opt offline.ExportOpt) (<-chan types.ExportRecord, <-chan error) {
	panic("implement me")
}

func (db *DB) Import(ctx context.Context, opt offline.ImportOpt) (int64, error) {
	panic("implement me")
}

func (db *DB) TableSchema(ctx context.Context, tableName string) (*types.DataTableSchema, error) {
	panic("implement me")
}

func (db *DB) TypeTag(dbType string) (string, error) {
	panic("implement me")
}
