package dbutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

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
	case types.POSTGRES, types.SNOWFLAKE:
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

func Quote(quote string, fields ...string) string {
	var rs []string
	for _, f := range fields {
		rs = append(rs, quote+f+quote)
	}
	return strings.Join(rs, ",")
}

func GetColumnFormat(backendType types.BackendType) (string, error) {
	var columnFormat string
	switch backendType {
	case types.POSTGRES, types.SNOWFLAKE:
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
	case types.POSTGRES, types.SNOWFLAKE:
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
