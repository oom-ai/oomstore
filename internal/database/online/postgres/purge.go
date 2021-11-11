package postgres

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

func (db *DB) Purge(ctx context.Context, revision *typesv2.Revision) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS "%s";`, getOnlineBatchTableName(revision.ID))
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
