package bigquery

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) Import(ctx context.Context, opt offline.ImportOpt) (int64, error) {
	tableID := opt.SnapshotTableName

	// Step 1: define table schema
	schema := make(bigquery.Schema, 0, len(opt.Features)+1)
	for _, h := range opt.Header {
		if h == opt.Entity.Name {
			schema = append(schema, &bigquery.FieldSchema{
				Name:        opt.Entity.Name,
				Type:        bigquery.StringFieldType,
				Description: "entity key",
			})
		} else {
			feature := opt.Features.Find(func(f *types.Feature) bool {
				return f.Name == h
			})
			if feature == nil {
				return 0, errdefs.Errorf("missing feature %s", h)
			}
			fieldType, err := convertValueTypeToBigQueryType(feature.ValueType)
			if err != nil {
				return 0, err
			}
			schema = append(schema, &bigquery.FieldSchema{
				Name:        feature.Name,
				Type:        fieldType,
				Description: feature.Description,
			})
		}
	}

	// Step 2: create offline table
	metaData := &bigquery.TableMetadata{
		Name:   opt.SnapshotTableName,
		Schema: schema,
	}
	tableRef := db.Dataset(db.datasetID).Table(tableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		return 0, errdefs.WithStack(err)
	}

	// Step 3: load data from source
	source := bigquery.NewReaderSource(opt.Source.Reader)
	source.Schema = schema
	loader := db.Dataset(db.datasetID).Table(tableID).LoaderFrom(source)
	job, err := loader.Run(ctx)
	if err != nil {
		return 0, errdefs.WithStack(err)
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return 0, errdefs.WithStack(err)
	}
	if err := status.Err(); err != nil {
		return 0, errdefs.WithStack(err)
	}
	return time.Now().UnixMilli(), nil
}

func convertValueTypeToBigQueryType(t types.ValueType) (bigquery.FieldType, error) {
	switch t {
	case types.String:
		return bigquery.StringFieldType, nil
	case types.Int64:
		return bigquery.IntegerFieldType, nil
	case types.Bool:
		return bigquery.BooleanFieldType, nil
	case types.Float64:
		return bigquery.FloatFieldType, nil
	case types.Bytes:
		return bigquery.BytesFieldType, nil
	case types.Time:
		return bigquery.TimeFieldType, nil
	default:
		return "", errdefs.Errorf("unsupported value type %s", t)
	}
}
