package cassandra

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/dbutil"

	"github.com/oom-ai/oomstore/pkg/errdefs"
)

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	err := db.Query(fmt.Sprintf("DROP TABLE IF EXISTS %s", dbutil.OnlineBatchTableName(revisionID))).
		WithContext(ctx).Exec()
	return errdefs.WithStack(err)
}
