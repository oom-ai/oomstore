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

func InsertRecordsToTable(ctx context.Context, dbOpt DBOpt, tableName string, records []interface{}, columns []string, backed types.BackendType) error {
	query, args, err := dbOpt.BuildInsertQuery(tableName, records, columns, backed)
	if err != nil {
		return err
	}
	if query == "" {
		return nil
	}

	return dbOpt.ExecContext(ctx, query, args)
}

func InsertRecordsToTableTx(tx *sqlx.Tx, ctx context.Context, tableName string, records []interface{}, columns []string, backendType types.BackendType) error {
	dbOpt := DBOpt{
		Backend: backendType,
	}
	query, args, err := dbOpt.BuildInsertQuery(tableName, records, columns, backendType)
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

func GetColumnFormat(backendType types.BackendType) (string, error) {
	var columnFormat string
	switch backendType {
	case types.BackendPostgres, types.BackendSnowflake, types.BackendRedshift:
		columnFormat = `"%s" %s`
	case types.BackendMySQL, types.BackendSQLite, types.BackendBigQuery:
		columnFormat = "`%s` %s"
	default:
		return "", fmt.Errorf("unsupported backend type %s", backendType)
	}
	return columnFormat, nil
}

func QuoteFn(backendType types.BackendType) func(...string) string {
	var q string
	switch backendType {
	case types.BackendPostgres, types.BackendSnowflake, types.BackendRedshift:
		q = `"`
	case types.BackendMySQL, types.BackendSQLite, types.BackendBigQuery:
		q = "`"
	default:
		panic(fmt.Sprintf("unsupported backend type %s", backendType))
	}
	return func(fields ...string) string {
		var rs []string
		for _, f := range fields {
			rs = append(rs, q+f+q)
		}
		return strings.Join(rs, ",")
	}
}
