package sqlutil

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func Purge(ctx context.Context, db *sqlx.DB, revisionID int, backend types.BackendType) error {
	qt := dbutil.QuoteFn(backend)
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, qt(dbutil.OnlineBatchTableName(revisionID)))
	_, err := db.ExecContext(ctx, query)
	return errdefs.WithStack(err)
}

func PurgeTx(ctx context.Context, tx *sqlx.Tx, tableName string, backend types.BackendType) error {
	qt := dbutil.QuoteFn(backend)
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, qt(tableName))
	_, err := tx.ExecContext(ctx, query)
	return errdefs.WithStack(err)
}
