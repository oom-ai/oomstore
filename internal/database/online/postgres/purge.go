package postgres

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) Purge(ctx context.Context, revision *types.Revision) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS "%s";`, getOnlineBatchTableName(revision.ID))
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
