package dbutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cast"

	"cloud.google.com/go/bigquery"
	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type DBOpt struct {
	Backend types.BackendType

	// Sqlx
	SqlxDB *sqlx.DB

	// BigQuery
	BigQueryDB *bigquery.Client
	DatasetID  *string
}

func (d *DBOpt) ExecContext(ctx context.Context, query string, args []interface{}) error {
	switch d.Backend {
	case types.BackendBigQuery:
		_, err := d.BigQueryDB.Query(query).Read(ctx)
		return err
	default:
		_, err := d.SqlxDB.ExecContext(ctx, d.SqlxDB.Rebind(query), args...)
		return err
	}
}

func (d *DBOpt) BuildInsertQuery(tableName string, records []interface{}, columns []string) (string, []interface{}, error) {
	if len(records) == 0 {
		return "", nil, nil
	}

	var columnStr string
	switch d.Backend {
	case types.BackendPostgres, types.BackendSnowflake, types.BackendRedshift:
		columnStr = Quote(`"`, columns...)
		tableName = fmt.Sprintf(`"%s"`, tableName)
	case types.BackendMySQL, types.BackendSQLite, types.BackendBigQuery:
		columnStr = Quote("`", columns...)
		tableName = fmt.Sprintf("`%s`", tableName)
	}

	switch d.Backend {
	case types.BackendBigQuery:
		values := make([]string, 0, len(records))
		for _, row := range records {
			values = append(values, fmt.Sprintf("(%s)", strings.Join(cast.ToStringSlice(row), ",")))
		}
		return fmt.Sprintf(`INSERT INTO %s.%s (%s) VALUES %s`, *d.DatasetID, tableName, columnStr, strings.Join(values, ",")), nil, nil
	default:
		valueFlags := make([]string, 0, len(records))
		for i := 0; i < len(records); i++ {
			valueFlags = append(valueFlags, "(?)")
		}
		return sqlx.In(
			fmt.Sprintf(`INSERT INTO %s (%s) VALUES %s`, tableName, columnStr, strings.Join(valueFlags, ",")),
			records...)
	}
}
