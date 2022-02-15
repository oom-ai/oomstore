package sqlutil

import (
	"context"
	"fmt"
	"time"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func AddTemporaryTableRecord(ctx context.Context, dbOpt dbutil.DBOpt, tableName string) error {
	if err := createTemporaryTableRecordTable(ctx, dbOpt); err != nil {
		return err
	}
	query := fmt.Sprintf(`INSERT INTO %s(table_name, create_time) VALUES(?,?)`, offline.TemporaryTableRecordTable)
	return dbOpt.ExecContext(ctx, query, []interface{}{tableName, time.Now().UnixMilli()})
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
`, offline.TemporaryTableRecordTable, tableNameDBType, createTimeDBType)
	return dbOpt.ExecContext(ctx, query, nil)
}

func DropTemporaryTables(ctx context.Context, db dbutil.DBOpt, tableNames []string) error {
	var err error
	for _, tableName := range tableNames {
		if tmpErr := dropTemporaryTable(ctx, db, tableName); tmpErr != nil {
			err = tmpErr
		}
	}
	return err
}

func dropTemporaryTable(ctx context.Context, db dbutil.DBOpt, tableName string) error {
	qt := dbutil.QuoteFn(db.Backend)
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s`, qt(tableName))
	if err := db.ExecContext(ctx, query, nil); err != nil {
		return errdefs.WithStack(err)
	}
	query = fmt.Sprintf(`DELETE FROM %s where table_name = ?`, offline.TemporaryTableRecordTable)
	if err := db.ExecContext(ctx, query, []interface{}{tableName}); err != nil {
		return errdefs.WithStack(err)
	}
	return nil
}
