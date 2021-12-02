package mysql

import (
	"context"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateEntityOpt) (int, error) {
	query := "INSERT INTO entity(name, length, description) VALUES(?, ?, ?)"
	res, err := sqlxCtx.ExecContext(ctx, query, opt.EntityName, opt.Length, opt.Description)
	if er, ok := err.(*mysql.MySQLError); ok {
		if er.Number == ER_DUP_ENTRY {
			return 0, fmt.Errorf("entity %s already exists", opt.EntityName)
		}
	}
	entityID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(entityID), err
}
