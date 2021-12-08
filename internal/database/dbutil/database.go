package dbutil

import (
	"fmt"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

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
