package bigquery

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/spf13/cast"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"
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

func (db *DB) DropTemporaryTable(ctx context.Context, tableNames []string) error {
	dbOpt := dbutil.DBOpt{Backend: Backend, BigQueryDB: db.Client, DatasetID: &db.datasetID}
	return sqlutil.DropTemporaryTables(ctx, dbOpt, tableNames)
}

func (db *DB) GetTemporaryTables(ctx context.Context, unixMilli int64) ([]string, error) {
	qt := dbutil.QuoteFn(Backend)
	query := fmt.Sprintf("SELECT table_name FROM %s.%s WHERE create_time < %s",
		db.datasetID, qt(offline.TemporaryTableRecordTable), cast.ToString(unixMilli))

	rows, err := db.Query(query).Read(ctx)
	if err != nil {
		tableNotFound, notFoundErr := dbutil.IsTableNotFoundError(err, Backend)
		if notFoundErr != nil {
			return nil, notFoundErr
		}
		if tableNotFound {
			return nil, nil
		}
		return nil, err
	}

	var tableNames []string
	for {
		recordMap := make(map[string]bigquery.Value)
		err = rows.Next(&recordMap)
		if err == iterator.Done {
			break
		}
		tableNames = append(tableNames, recordMap["table_name"].(string))
	}
	return tableNames, nil
}
