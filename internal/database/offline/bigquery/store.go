package bigquery

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"

	"github.com/spf13/cast"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/bigquery"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"google.golang.org/api/option"
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

func (db *DB) TableSchema(ctx context.Context, tableName string) (*types.DataTableSchema, error) {
	q := fmt.Sprintf(`SELECT column_name, data_type FROM %s.INFORMATION_SCHEMA.COLUMNS WHERE table_name = "%s"`, db.datasetID, tableName)
	rows, err := db.Query(q).Read(ctx)
	if err != nil {
		return nil, err
	}
	var schema types.DataTableSchema
	for {
		recordMap := make(map[string]bigquery.Value)
		err = rows.Next(&recordMap)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		valueType, err := dbutil.ValueType(Backend, cast.ToString(recordMap["data_type"]))
		if err != nil {
			return nil, err
		}
		schema.Fields = append(schema.Fields, types.DataTableFieldSchema{
			Name:      cast.ToString(recordMap["column_name"]),
			ValueType: valueType,
		})
	}

	return &schema, nil
}

func (db *DB) Snapshot(ctx context.Context, opt offline.SnapshotOpt) error {
	dbOpt := dbutil.DBOpt{
		Backend:    Backend,
		BigQueryDB: db.Client,
		DatasetID:  &db.datasetID,
	}
	return sqlutil.Snapshot(ctx, dbOpt, opt)
}

func (db *DB) CreateTable(ctx context.Context, opt offline.CreateTableOpt) error {
	schema := dbutil.BuildTableSchema(opt.TableName, opt.Entity, opt.IsCDC, opt.Features, nil, Backend)
	if _, err := db.Client.Query(schema).Run(ctx); err != nil {
		return err
	}
	return nil
}
