package cassandra

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
)

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	err := db.Query(fmt.Sprintf("DROP TABLE IF EXISTS %s", sqlutil.OnlineBatchTableName(revisionID))).
		WithContext(ctx).Exec()
	return errors.WithStack(err)
}
