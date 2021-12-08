package dbutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type RowMap = map[string]interface{}

func BuildConditions(equal map[string]interface{}, in map[string]interface{}) ([]string, []interface{}, error) {
	cond := make([]string, 0)
	args := make([]interface{}, 0)
	for key, value := range equal {
		cond = append(cond, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}
	for key, value := range in {
		s, inArgs, err := sqlx.In(fmt.Sprintf("%s IN (?)", key), value)
		if err != nil {
			return nil, nil, err
		}
		cond = append(cond, s)
		args = append(args, inArgs...)
	}
	return cond, args, nil
}

func Quote(quote string, fields ...string) string {
	var rs []string
	for _, f := range fields {
		rs = append(rs, quote+f+quote)
	}
	return strings.Join(rs, ",")
}

func TempTable(prefix string) string {
	return fmt.Sprintf("tmp_%s_%d", prefix, time.Now().UnixNano())
}

func buildQueryAndArgsForInsertRecords(tableName string, records []interface{}, columns []string, backendType types.BackendType) (string, []interface{}, error) {
	if len(records) == 0 {
		return "", nil, nil
	}
	valueFlags := make([]string, 0, len(records))
	for i := 0; i < len(records); i++ {
		valueFlags = append(valueFlags, "(?)")
	}

	var columnStr string
	switch backendType {
	case types.POSTGRES:
		columnStr = Quote(`"`, columns...)
		tableName = fmt.Sprintf(`"%s"`, tableName)
	case types.MYSQL:
		columnStr = Quote("`", columns...)
		tableName = fmt.Sprintf("`%s`", tableName)
	}
	return sqlx.In(
		fmt.Sprintf(`INSERT INTO %s (%s) VALUES %s`, tableName, columnStr, strings.Join(valueFlags, ",")),
		records...)
}

func InsertRecordsToTable(db *sqlx.DB, ctx context.Context, tableName string, records []interface{}, columns []string, backendType types.BackendType) error {
	query, args, err := buildQueryAndArgsForInsertRecords(tableName, records, columns, backendType)
	if err != nil {
		return err
	}
	if query == "" {
		return nil
	}
	if _, err := db.ExecContext(ctx, db.Rebind(query), args...); err != nil {
		return err
	}
	return nil
}

func InsertRecordsToTableTx(tx *sqlx.Tx, ctx context.Context, tableName string, records []interface{}, columns []string, backendType types.BackendType) error {
	query, args, err := buildQueryAndArgsForInsertRecords(tableName, records, columns, backendType)
	if err != nil {
		return err
	}
	if query == "" {
		return nil
	}
	if _, err := tx.ExecContext(ctx, tx.Rebind(query), args...); err != nil {
		return err
	}
	return nil
}

// Build MySQL data source name
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

func GetColumnFormat(backendType types.BackendType) (string, error) {
	var columnFormat string
	switch backendType {
	case types.POSTGRES:
		columnFormat = `"%s" %s`
	case types.MYSQL:
		columnFormat = "`%s` %s"
	default:
		return "", fmt.Errorf("unsupported backend type %s", backendType)
	}
	return columnFormat, nil
}

func QuoteFn(backendType types.BackendType) (func(...string) string, error) {
	var quote string
	switch backendType {
	case types.POSTGRES:
		quote = `"`
	case types.MYSQL:
		quote = "`"
	default:
		return nil, fmt.Errorf("unsupported backend type %s", backendType)
	}
	return func(fields ...string) string {
		return Quote(quote, fields...)
	}, nil
}
