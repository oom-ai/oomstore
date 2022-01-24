package sqlutil

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/offline"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// SqlxTableSchema returns the schema of the given table.
func SqlxTableSchema(ctx context.Context, db *sqlx.DB, backend types.BackendType, rows *sqlx.Rows, opt offline.TableSchemaOpt) (*types.DataTableSchema, error) {
	defer rows.Close()

	var schema types.DataTableSchema
	for rows.Next() {
		var fieldName, dbValueType string
		if err := rows.Scan(&fieldName, &dbValueType); err != nil {
			return nil, errdefs.WithStack(err)
		}
		valueType, err := dbutil.ValueType(backend, dbValueType)
		if err != nil {
			return nil, err
		}
		schema.Fields = append(schema.Fields, types.DataTableFieldSchema{
			Name:      fieldName,
			ValueType: valueType,
		})
	}
	if len(schema.Fields) == 0 {
		return nil, errdefs.Errorf("table not found")
	}
	if opt.CheckTimeRange {
		timeRange, err := getCdcTimeRange(ctx, db, opt.TableName, backend)
		if err != nil {
			return nil, err
		}
		schema.TimeRange = *timeRange
	}
	return &schema, nil
}

func getCdcTimeRange(ctx context.Context, db *sqlx.DB, tableName string, backend types.BackendType) (*types.DataTableTimeRange, error) {
	qt := dbutil.QuoteFn(backend)
	var timeRange types.DataTableTimeRange
	query := fmt.Sprintf(`
		SELECT
			MIN(%s) AS %s,
			MAX(%s) AS %s
		FROM %s`, qt("unix_milli"), qt("min_unix_milli"), qt("unix_milli"), qt("max_unix_milli"), qt(tableName))

	if err := db.GetContext(ctx, &timeRange, query); err != nil {
		return nil, errdefs.WithStack(err)
	}
	return &timeRange, nil
}
