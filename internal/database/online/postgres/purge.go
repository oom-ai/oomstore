package postgres

import (
	"context"
	"fmt"
)

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS "%s";`, getOnlineBatchTableName(revisionID))
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
