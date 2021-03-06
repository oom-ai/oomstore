package mysql

import (
	"context"

	"github.com/go-sql-driver/mysql"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/errdefs"
)

func createEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateEntityOpt) (int, error) {
	query := "INSERT INTO entity(name, description) VALUES(?, ?)"
	res, err := sqlxCtx.ExecContext(ctx, query, opt.EntityName, opt.Description)
	if err != nil {
		if er, ok := err.(*mysql.MySQLError); ok {
			if er.Number == ER_DUP_ENTRY {
				return 0, errdefs.Errorf("entity %s already exists", opt.EntityName)
			}
		}
		return 0, errdefs.WithStack(err)
	}

	entityID, err := res.LastInsertId()
	if err != nil {
		return 0, errdefs.WithStack(err)
	}
	return int(entityID), nil
}
