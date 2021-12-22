package bigquery

import (
	"context"
	"fmt"

	"github.com/spf13/cast"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/bigquery"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"google.golang.org/api/option"
)

var BigQueryTypeMap = map[string]types.ValueType{
	"bool":     types.BOOL,
	"bytes":    types.BYTES,
	"datetime": types.TIME,
	"string":   types.STRING,

	"bigint":   types.INT64,
	"smallint": types.INT64,
	"int64":    types.INT64,
	"integer":  types.INT64,
	"int":      types.INT64,

	"float64": types.FLOAT64,
	"numeric": types.FLOAT64,
	"decimal": types.FLOAT64,
}

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
		valueType, err := db.TypeTag(cast.ToString(recordMap["data_type"]))
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

func (db *DB) TypeTag(dbType string) (types.ValueType, error) {
	return sqlutil.GetValueType(BigQueryTypeMap, dbType)
}
