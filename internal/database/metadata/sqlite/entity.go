package sqlite

import (
	"context"

	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateEntityOpt) (int, error) {
	query := "INSERT INTO entity(name, description) VALUES(?, ?)"
	res, err := sqlxCtx.ExecContext(ctx, query, opt.EntityName, opt.Description)
	if err != nil {
		if er, ok := err.(sqlite3.Error); ok {
			if er.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, errors.Errorf("entity %s already exists", opt.EntityName)
			}
		}
		return 0, err
	}

	entityID, err := res.LastInsertId()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return int(entityID), nil
}
