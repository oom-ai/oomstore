package cassandra

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
)

func (db *DB) PrepareStreamTable(ctx context.Context, opt online.PrepareStreamTableOpt) error {
	tableName := sqlutil.OnlineStreamTableName(opt.GroupID)

	if opt.Feature == nil {
		schema, err := sqlutil.CreateStreamTableSchema(ctx, tableName, opt.Entity, Backend)
		if err != nil {
			return err
		}

		return db.Query(schema).WithContext(ctx).Exec()
	}

	dbValueType, err := dbutil.DBValueType(Backend, opt.Feature.ValueType)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("ALTER TABLE %s ADD %s %s", tableName, opt.Feature.Name, dbValueType)

	return db.Query(sql).WithContext(ctx).Exec()
}

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	tableName := sqlutil.OnlineStreamTableName(opt.GroupID)

	cond := sqlutil.BuildPushCondition(opt, Backend)
	// cassandra's `insert` is equivalent to `insert_or_update`.
	// see: https://cassandra.apache.org/doc/latest/cassandra/cql/dml.html#insert-statement
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)",
		tableName,
		cond.Inserts,
		cond.InsertPlaceholders,
	)

	return db.Query(query, cond.InsertValues...).WithContext(ctx).Exec()
}
