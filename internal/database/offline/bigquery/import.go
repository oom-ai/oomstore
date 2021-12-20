package bigquery

import (
	"context"
	"fmt"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/oom-ai/oomstore/internal/database/offline"
)

func (db *DB) Import(ctx context.Context, opt offline.ImportOpt) (int64, error) {
	tableID := opt.DataTableName

	// Step 1: define table schema
	schema := make(bigquery.Schema, 0, len(opt.Features)+1)
	schema = append(schema, &bigquery.FieldSchema{
		Name:        opt.Entity.Name,
		Type:        bigquery.StringFieldType,
		Description: "entity key",
	})
	for _, f := range opt.Features {
		fieldType, err := convertValueTypeToBigQueryType(f.ValueType)
		if err != nil {
			return 0, err
		}
		schema = append(schema, &bigquery.FieldSchema{
			Name:        f.Name,
			Type:        fieldType,
			Description: f.Description,
		})
	}

	// Step 2: create offline etable
	metaData := &bigquery.TableMetadata{
		Name:   opt.DataTableName,
		Schema: schema,
	}
	tableRef := db.Dataset(db.datasetID).Table(tableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		return 0, err
	}

	// Step 3: load data from source
	source := bigquery.NewReaderSource(opt.Source.Reader)
	source.Schema = schema
	loader := db.Dataset(db.datasetID).Table(tableID).LoaderFrom(source)
	job, err := loader.Run(ctx)
	if err != nil {
		return 0, err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return 0, err
	}
	if err := status.Err(); err != nil {
		return 0, err
	}
	return time.Now().UnixMilli(), nil
}

func convertValueTypeToBigQueryType(t string) (bigquery.FieldType, error) {
	switch t {
	case types.STRING:
		return bigquery.StringFieldType, nil
	case types.INT64:
		return bigquery.IntegerFieldType, nil
	case types.BOOL:
		return bigquery.BooleanFieldType, nil
	case types.FLOAT64:
		return bigquery.FloatFieldType, nil
	case types.BYTES:
		return bigquery.BytesFieldType, nil
	case types.TIME:
		return bigquery.TimeFieldType, nil
	default:
		return "", fmt.Errorf("unsupported value type %s", t)
	}
}
