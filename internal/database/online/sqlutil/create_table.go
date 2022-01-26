package sqlutil

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/online"
)

func CreateTable(ctx context.Context, dbOpt dbutil.DBOpt, opt online.CreateTableOpt) error {
	// Step 1: drop existing table
	if err := dbutil.DropTable(ctx, dbOpt, opt.TableName); err != nil {
		return err
	}
	// Step 2: create new table
	schema := dbutil.BuildTableSchema(opt.TableName, opt.EntityName, false, opt.Features, []string{opt.EntityName}, dbOpt.Backend)
	err := dbOpt.ExecContext(ctx, schema, nil)
	if err != nil {
		return errdefs.WithStack(err)
	}
	return nil
}
