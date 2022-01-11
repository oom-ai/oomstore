package sqlite

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateEntityOpt) (int, error) {
	query := "INSERT INTO entity(name, description) VALUES(?, ?)"
	res, err := sqlxCtx.ExecContext(ctx, query, opt.EntityName, opt.Description)
	if err != nil {
		if er, ok := err.(*sqlite.Error); ok {
			if er.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return 0, errdefs.Errorf("entity %s already exists", opt.EntityName)
			}
		}

		return 0, err
	}

	entityID, err := res.LastInsertId()
	if err != nil {
		return 0, errdefs.WithStack(err)
	}
	return int(entityID), nil
}
