package sqlutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func AddTemporaryTableRecord(ctx context.Context, dbOpt dbutil.DBOpt, tableName string) error {
	if err := createTemporaryTableRecordTable(ctx, dbOpt); err != nil {
		return err
	}
	unQt := dbutil.UnQuoteFn(dbOpt.Backend)
	tableName = unQt(tableName)
	if dbOpt.Backend == types.BackendBigQuery {
		tableName = fmt.Sprintf(`"%s"`, tableName)
	}
	query := fmt.Sprintf(`INSERT INTO %s (table_name, create_time) VALUES(?,?)`, buildTableName(dbOpt, offline.TemporaryTableRecordTable))
	return dbOpt.ExecContext(ctx, query, tableName, time.Now().UnixMilli())
}

func createTemporaryTableRecordTable(ctx context.Context, dbOpt dbutil.DBOpt) error {
	tableNameDBType, err := dbutil.DBValueType(dbOpt.Backend, types.String)
	if err != nil {
		return err
	}

	createTimeDBType, err := dbutil.DBValueType(dbOpt.Backend, types.Int64)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
	table_name %s,
	create_time %s
)
`, buildTableName(dbOpt, offline.TemporaryTableRecordTable), tableNameDBType, createTimeDBType)
	return dbOpt.ExecContext(ctx, query)
}

func GetTemporaryTables(ctx context.Context, db *sqlx.DB, backend types.BackendType, unixMill int64) ([]string, error) {
	var tableName string
	if backend == types.BackendSnowflake {
		tableName = fmt.Sprintf(`PUBLIC."%s"`, offline.TemporaryTableRecordTable)
	} else {
		tableName = offline.TemporaryTableRecordTable
	}
	query := fmt.Sprintf("SELECT table_name FROM %s WHERE create_time < ?", tableName)

	rows, err := db.QueryContext(ctx, db.Rebind(query), unixMill)
	if err != nil {
		tableNotFound, notFoundErr := dbutil.IsTableNotFoundError(err, backend)
		if notFoundErr != nil {
			return nil, notFoundErr
		}
		if tableNotFound {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tableNames = append(tableNames, tableName)
	}
	return tableNames, nil
}

func DropTemporaryTables(ctx context.Context, db dbutil.DBOpt, tableNames []string) error {
	for _, tableName := range tableNames {
		query := fmt.Sprintf(`DROP TABLE IF EXISTS %s`, tableName)
		if err := db.ExecContext(ctx, query); err != nil {
			return err
		}
	}

	unQt := dbutil.UnQuoteFn(db.Backend)
	for i := 0; i < len(tableNames); i++ {
		tableNames[i] = unQt(tableNames[i])
		if db.Backend == types.BackendBigQuery {
			tableNames[i] = fmt.Sprintf(`"%s"`, tableNames[i])
		}
	}

	cond, args, err := dbutil.BuildConditions(nil, map[string]interface{}{
		"table_name": tableNames,
	})
	if err != nil {
		return nil
	}
	if len(cond) > 0 {
		query := fmt.Sprintf("DELETE FROM %s WHERE %s",
			buildTableName(db, offline.TemporaryTableRecordTable),
			strings.Join(cond, " AND "))
		return db.ExecContext(ctx, query, args...)
	}
	return nil
}
