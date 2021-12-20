package cassandra

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
)

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	return db.Query(fmt.Sprintf("DROP TABLE IF EXISTS %s", sqlutil.OnlineTableName(revisionID))).
		WithContext(ctx).Exec()
}
