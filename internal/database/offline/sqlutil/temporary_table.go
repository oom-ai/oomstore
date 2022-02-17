package sqlutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func AddTemporaryTableRecord(ctx context.Context, dbOpt dbutil.DBOpt, tableName string) error {
	if err := createTemporaryTableRecordTable(ctx, dbOpt); err != nil {
		// The logic of the temporary table should not affect the main process, so nil is returned here.
		// TODO: Print log in the cloud service version of oomstore
		return nil
	}
	unQt := dbutil.UnQuoteFn(dbOpt.Backend)
	tableName = unQt(tableName)
	if dbOpt.Backend == types.BackendBigQuery {
		tableName = fmt.Sprintf(`"%s"`, tableName)
	}
	query := fmt.Sprintf(`INSERT INTO %s (table_name, create_time) VALUES(?,?)`, buildTableName(dbOpt, offline.TemporaryTableRecordTable))
	if err := dbOpt.ExecContext(ctx, query, tableName, time.Now().UnixMilli()); err != nil {
		// The logic of the temporary table should not affect the main process, so nil is returned here.
		// TODO: Print log in the cloud service version of oomstore
		return nil
	}
	return nil
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

func DropTemporaryTables(ctx context.Context, db dbutil.DBOpt, tableNames []string) error {
	for _, tableName := range tableNames {
		query := fmt.Sprintf(`DROP TABLE IF EXISTS %s`, tableName)
		if err := db.ExecContext(ctx, query); err != nil {
			// The logic of the temporary table should not affect the main process, so nil is returned here.
			// TODO: Print log in the cloud service version of oomstore
			return nil
		}
	}

	unQt := dbutil.UnQuoteFn(db.Backend)
	for i := 0; i < len(tableNames); i++ {
		tableNames[i] = unQt(tableNames[i])
	}
	cond, args, err := dbutil.BuildConditions(nil, map[string]interface{}{
		"table_name": tableNames,
	})
	if err != nil {
		// The logic of the temporary table should not affect the main process, so nil is returned here.
		// TODO: Print log in the cloud service version of oomstore
		return nil
	}
	if len(cond) > 0 {
		query := fmt.Sprintf("DELETE FROM %s WHERE %s",
			buildTableName(db, offline.TemporaryTableRecordTable),
			strings.Join(cond, " AND "))
		if err := db.ExecContext(ctx, query, args...); err != nil {
			// The logic of the temporary table should not affect the main process, so nil is returned here.
			// TODO: Print log in the cloud service version of oomstore
			return nil
		}
	}
	return nil
}
