package sqlutil

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func Purge(ctx context.Context, db *sqlx.DB, revisionID int, backend types.BackendType) error {
	qt, err := dbutil.QuoteFn(backend)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, qt(OnlineBatchTableName(revisionID)))
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
