package sqlutil

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// SqlxTableSchema returns the schema of the given table.
func SqlxTableSchema(ctx context.Context, store offline.Store, backend types.BackendType, rows *sqlx.Rows) (*types.DataTableSchema, error) {
	defer rows.Close()

	var schema types.DataTableSchema
	for rows.Next() {
		var fieldName, dbValueType string
		if err := rows.Scan(&fieldName, &dbValueType); err != nil {
			return nil, err
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
		return nil, fmt.Errorf("table not found")
	}
	return &schema, nil
}
