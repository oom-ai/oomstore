package bigquery

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"google.golang.org/api/option"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	Backend = types.BackendBigQuery
)

var _ offline.Store = &DB{}

type DB struct {
	*bigquery.Client
	datasetID string
}

func Open(ctx context.Context, opt *types.BigQueryOpt) (*DB, error) {
	client, err := bigquery.NewClient(ctx, opt.ProjectID, option.WithCredentialsJSON([]byte(opt.Credentials)))
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	return &DB{
		Client:    client,
		datasetID: opt.DatasetID,
	}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	q := db.Client.Query("SELECT 1")
	_, err := q.Read(ctx)
	return errdefs.WithStack(err)
}

func (db *DB) Snapshot(ctx context.Context, opt offline.SnapshotOpt) error {
	dbOpt := dbutil.DBOpt{
		Backend:    Backend,
		BigQueryDB: db.Client,
		DatasetID:  &db.datasetID,
	}
	return errdefs.WithStack(sqlutil.Snapshot(ctx, dbOpt, opt))
}

func (db *DB) CreateTable(ctx context.Context, opt offline.CreateTableOpt) error {
	dbOpt := dbutil.DBOpt{Backend: Backend, BigQueryDB: db.Client, DatasetID: &db.datasetID}
	return sqlutil.CreateTable(ctx, dbOpt, opt)
}

func (db *DB) Push(ctx context.Context, opt offline.PushOpt) error {
	dbOpt := dbutil.DBOpt{Backend: Backend, BigQueryDB: db.Client, DatasetID: &db.datasetID}
	if err := sqlutil.Push(ctx, dbOpt, opt); err != nil {
		return err
	}
	return nil
}

func (db *DB) DropTemporaryTable(ctx context.Context, unixmilli int64) error {
	return sqlutil.DropTemporaryTables(ctx, dbutil.DBOpt{
		BigQueryDB: db.Client,
		Backend:    Backend,
	}, offline.DropTemporaryTableParams{
		UnixMilli: &unixmilli,
	})
}
