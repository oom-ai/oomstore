package bigquery

import (
	"context"
	"fmt"

	"github.com/spf13/cast"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/bigquery"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"google.golang.org/api/option"
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
		schema.Fields = append(schema.Fields, types.DataTableFieldSchema{
			Name:      cast.ToString(recordMap["column_name"]),
			ValueType: cast.ToString(recordMap["data_type"]),
		})
	}

	return &schema, nil
}

func (db *DB) TypeTag(dbType string) (string, error) {
	return TypeTag(dbType)
}
