package sqlutil

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func CreateTable(ctx context.Context, db *sqlx.DB, opt offline.CreateTableOpt, backend types.BackendType) error {
	schema := dbutil.BuildTableSchema(opt.TableName, opt.Entity, opt.WithUnixMillis, opt.Features, backend)
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return err
	}
	return nil
}
