package cassandra

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) PrepareStreamTable(ctx context.Context, opt online.PrepareStreamTableOpt) error {
	tableName := sqlutil.OnlineStreamTableName(opt.GroupID)

	if opt.Feature == nil {
		schema, err := sqlutil.CreateStreamTableSchema(ctx, tableName, opt.Entity, types.BackendMySQL)
		if err != nil {
			return err
		}

		return db.Query(schema).WithContext(ctx).Exec()
	}

	dbValueType, err := dbutil.DBValueType(types.BackendCassandra, opt.Feature.ValueType)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, opt.Feature.Name, dbValueType)
	return db.Query(sql).WithContext(ctx).Exec()
}