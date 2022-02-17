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
	schema := dbutil.BuildTableSchema(dbutil.BuildTableSchemaParams{
		TableName:    opt.TableName,
		EntityName:   opt.EntityName,
		HasUnixMilli: false,
		Features:     opt.Features,
		PrimaryKeys:  []string{opt.EntityName},
		Backend:      dbOpt.Backend,
	})
	err := dbOpt.ExecContext(ctx, schema)
	if err != nil {
		return errdefs.WithStack(err)
	}
	return nil
}
