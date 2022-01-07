package dbutil

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/pkg/errors"
	"github.com/snowflakedb/gosnowflake"
	"google.golang.org/api/googleapi"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func OpenSQLite(dbFile string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", dbFile)
	return db, errors.WithStack(err)
}

func OpenMysqlDB(host, port, user, password, database string) (*sqlx.DB, error) {
	cfg := mysql.NewConfig()
	cfg.Net = fmt.Sprintf("(%s:%s)", host, port)
	cfg.User = user
	cfg.Passwd = password
	cfg.DBName = database
	cfg.ParseTime = true

	db, err := sqlx.Open("mysql", cfg.FormatDSN())
	return db, errors.WithStack(err)
}

func OpenPostgresDB(host, port, user, password, database string) (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			password,
			host,
			port,
			database),
	)
	return db, errors.WithStack(err)
}

var OpenRedshiftDB = OpenPostgresDB

func DeserializeString(i interface{}, backend types.BackendType) string {
	switch backend {
	case types.BackendMySQL:
		return string(i.([]byte))
	default:
		return i.(string)
	}
}

// TODO: Should return an error when bakcend is not supported ?
func IsTableNotFoundError(err error, backend types.BackendType) bool {
	err = errors.Cause(err)
	switch backend {
	case types.BackendSQLite:
		if sqliteErr, ok := err.(*sqlite.Error); ok {
			return sqliteErr.Code() == sqlite3.SQLITE_CORE
		}

	// https://dev.mysql.com/doc/mysql-errors/5.7/en/server-error-reference.html#error_er_no_such_table
	case types.BackendMySQL:
		if e2, ok := err.(*mysql.MySQLError); ok {
			return e2.Number == 1146
		}
	case types.BackendPostgres, types.BackendRedshift:
		if e2, ok := err.(*pq.Error); ok {
			return e2.Code == pgerrcode.UndefinedTable
		}
	case types.BackendSnowflake:
		if e2, ok := err.(*gosnowflake.SnowflakeError); ok {
			return e2.Number == gosnowflake.ErrObjectNotExistOrAuthorized
		}
	// https://cloud.google.com/bigquery/docs/error-messages
	case types.BackendBigQuery:
		if e2, ok := err.(*googleapi.Error); ok {
			return e2.Code == 404
		}
	}
	return false
}
