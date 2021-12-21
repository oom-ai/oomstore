package dbutil

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func OpenSQLite(dbFile string) (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", dbFile)
}

func OpenMysqlDB(host, port, user, password, database string) (*sqlx.DB, error) {
	cfg := mysql.NewConfig()
	cfg.Net = fmt.Sprintf("(%s:%s)", host, port)
	cfg.User = user
	cfg.Passwd = password
	cfg.DBName = database
	cfg.ParseTime = true

	return sqlx.Open("mysql", cfg.FormatDSN())
}

func OpenPostgresDB(host, port, user, password, database string) (*sqlx.DB, error) {
	return sqlx.Open(
		"postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			password,
			host,
			port,
			database),
	)
}

func DeserializeString(i interface{}, backend types.BackendType) string {
	switch backend {
	case types.MYSQL:
		return string(i.([]byte))
	default:
		return i.(string)
	}
}

func IsTableNotFoundError(err error, backend types.BackendType) bool {
	switch backend {
	case types.SQLite:
		// https://github.com/mattn/go-sqlite3/issues/244
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			return sqliteErr.Code == sqlite3.ErrError
		}

	// https://dev.mysql.com/doc/mysql-errors/5.7/en/server-error-reference.html#error_er_no_such_table
	case types.MYSQL:
		if e2, ok := err.(*mysql.MySQLError); ok {
			return e2.Number == 1146
		}
	case types.POSTGRES:
		if e2, ok := err.(*pq.Error); ok {
			return e2.Code == pgerrcode.UndefinedTable
		}
	}
	return false
}
