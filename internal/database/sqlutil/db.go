package sqlutil

import (
	"context"
	"fmt"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
)

func Purge(ctx context.Context, db *sqlx.DB, revisionID int, backend types.BackendType) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, dbutil.QuoteByBackend(backend, OnlineTableName(revisionID)))
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

func OnlineTableName(revisionID int) string {
	return fmt.Sprintf("online_%d", revisionID)
}
