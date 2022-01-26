package cassandra

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"
)

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	tableName := dbutil.OnlineStreamTableName(opt.GroupID)

	cond := sqlutil.BuildPushCondition(opt, Backend)
	// cassandra's `insert` is equivalent to `insert_or_update`.
	// see: https://cassandra.apache.org/doc/latest/cassandra/cql/dml.html#insert-statement
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)",
		tableName,
		cond.Inserts,
		cond.InsertPlaceholders,
	)

	err := db.Query(query, cond.InsertValues...).WithContext(ctx).Exec()
	return errdefs.WithStack(err)
}
