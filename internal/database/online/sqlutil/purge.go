package sqlutil

import (
	"context"
	"fmt"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
)

func Purge(ctx context.Context, db *sqlx.DB, revisionID int, backend types.BackendType) error {
	qt, err := dbutil.QuoteFn(backend)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, qt(OnlineTableName(revisionID)))
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
