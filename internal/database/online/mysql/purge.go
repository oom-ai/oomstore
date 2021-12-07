package mysql

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/online/sqlutil"
)

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	return sqlutil.Purge(ctx, db.DB, revisionID, BackendType)
}
