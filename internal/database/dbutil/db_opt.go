package dbutil

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/gocql/gocql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type DBOpt struct {
	Backend types.BackendType

	// Sqlx
	SqlxDB *sqlx.DB

	// BigQuery
	BigQueryDB *bigquery.Client
	DatasetID  *string

	// Cassandra
	CassandraDB *gocql.Session
}

func (d *DBOpt) ExecContext(ctx context.Context, query string, args []interface{}) error {
	switch d.Backend {
	case types.BackendBigQuery:
		for _, arg := range args {
			query = strings.Replace(query, "?", cast.ToString(arg), 1)
		}
		_, err := d.BigQueryDB.Query(query).Read(ctx)
		return errdefs.WithStack(err)
	case types.BackendCassandra:
		return errdefs.WithStack(d.CassandraDB.Query(query).Exec())
	default:
		_, err := d.SqlxDB.ExecContext(ctx, d.SqlxDB.Rebind(query), args...)
		return errdefs.WithStack(err)
	}
}

func (d *DBOpt) BuildInsertQuery(tableName string, records []interface{}, columns []string) (string, []interface{}, error) {
	if len(records) == 0 {
		return "", nil, nil
	}

	qt := QuoteFn(d.Backend)
	columnStr := qt(columns...)
	tableName = qt(tableName)

	valueFlags := make([]string, 0, len(records))
	for i := 0; i < len(records); i++ {
		valueFlags = append(valueFlags, "(?)")
	}
	if d.Backend == types.BackendBigQuery {
		tableName = fmt.Sprintf("%s.%s", *d.DatasetID, tableName)
	}
	query, args, err := sqlx.In(
		fmt.Sprintf(`INSERT INTO %s (%s) VALUES %s`, tableName, columnStr, strings.Join(valueFlags, ",")),
		records...)
	return query, args, errdefs.WithStack(err)
}
