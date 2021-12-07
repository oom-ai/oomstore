package mysql

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/sqlutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	return sqlutil.Purge(ctx, db.DB, revisionID, types.MYSQL)
}
