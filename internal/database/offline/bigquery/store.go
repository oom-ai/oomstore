package bigquery

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var _ offline.Store = &DB{}

type DB struct {
	*bigquery.Client
	datasetID string
}

func Open(ctx context.Context, opt *types.BigQueryOpt) (*DB, error) {
	client, err := bigquery.NewClient(ctx, opt.ProjectID)
	if err != nil {
		return nil, err
	}
	return &DB{
		Client:    client,
		datasetID: opt.DatasetID,
	}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	q := db.Client.Query("SELECT 1")
	_, err := q.Read(ctx)
	return err
}

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	panic("implement me")
}

func (db *DB) TableSchema(ctx context.Context, tableName string) (*types.DataTableSchema, error) {
	panic("implement me")
}

func (db *DB) TypeTag(dbType string) (string, error) {
	panic("implement me")
}
