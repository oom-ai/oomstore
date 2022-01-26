package bigquery

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
	"google.golang.org/api/iterator"
)

func (db *DB) TableSchema(ctx context.Context, opt offline.TableSchemaOpt) (*types.DataTableSchema, error) {
	q := fmt.Sprintf(`SELECT column_name, data_type FROM %s.INFORMATION_SCHEMA.COLUMNS WHERE table_name = "%s"`, db.datasetID, opt.TableName)
	rows, err := db.Query(q).Read(ctx)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	var schema types.DataTableSchema
	for {
		recordMap := make(map[string]bigquery.Value)
		err = rows.Next(&recordMap)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errdefs.WithStack(err)
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
	if opt.CheckTimeRange {
		dbOpt := dbutil.DBOpt{
			Backend:    Backend,
			BigQueryDB: db.Client,
			DatasetID:  &db.datasetID,
		}
		timeRange, err := getTableTimeRange(ctx, dbOpt, opt.TableName)
		if err != nil {
			return nil, err
		}
		schema.TimeRange = *timeRange
	}
	return &schema, nil
}

func getTableTimeRange(ctx context.Context, dbOpt dbutil.DBOpt, tableName string) (*types.DataTableTimeRange, error) {
	q := fmt.Sprintf(`
		SELECT
			MIN(unix_milli) AS min_unix_milli,
			MAX(unix_milli) AS max_unix_milli
		FROM %s.%s`, *dbOpt.DatasetID, tableName)
	rows, err := dbOpt.BigQueryDB.Query(q).Read(ctx)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	recordMap := make(map[string]bigquery.Value)
	err = rows.Next(&recordMap)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	min := cast.ToInt64(recordMap["min_unix_milli"])
	max := cast.ToInt64(recordMap["max_unix_milli"])

	return &types.DataTableTimeRange{
		MinUnixMilli: &min,
		MaxUnixMilli: &max,
	}, nil
}
