package sqlutil

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func CreateTable(ctx context.Context, db *sqlx.DB, opt online.CreateTableOpt, backend types.BackendType) error {
	// Step 1: drop existing table
	if err := dbutil.DropTable(ctx, db, opt.TableName); err != nil {
		return err
	}
	// Step 2: create new table
	schema := dbutil.BuildTableSchema(opt.TableName, opt.EntityName, false, opt.Features, []string{opt.EntityName}, backend)
	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return errdefs.WithStack(err)
	}
	return nil
}
