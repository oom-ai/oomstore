package sqlite

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func CreateDatabase(ctx context.Context, opt *types.SQLiteOpt) (err error) {
	db, err := dbutil.OpenSQLite(opt.DBFile)
	if err != nil {
		return err
	}
	defer db.Close()
	return createMetaSchemas(ctx, db)
}

func createMetaSchemas(ctx context.Context, db *sqlx.DB) (err error) {
	return dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		for _, schema := range META_TABLE_SCHEMAS {
			if _, err = tx.ExecContext(ctx, schema); err != nil {
				return err
			}
		}

		for _, schema := range META_VIEW_SCHEMAS {
			if _, err = tx.ExecContext(ctx, schema); err != nil {
				return err
			}
		}

		for table := range META_TABLE_SCHEMAS {
			trigger := strings.ReplaceAll(TRIGGER_TEMPLATE, `{{TABLE_NAME}}`, table)
			if _, err = tx.ExecContext(ctx, trigger); err != nil {
				return err
			}
		}

		return nil
	})
}
