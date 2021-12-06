package mysql

import (
	"context"
	"fmt"

	"github.com/ethhte88/oomstore/internal/database/metadata"
	"github.com/go-sql-driver/mysql"
)

func createEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateEntityOpt) (int, error) {
	query := "INSERT INTO entity(name, length, description) VALUES(?, ?, ?)"
	res, err := sqlxCtx.ExecContext(ctx, query, opt.EntityName, opt.Length, opt.Description)
	if err != nil {
		if er, ok := err.(*mysql.MySQLError); ok {
			if er.Number == ER_DUP_ENTRY {
				return 0, fmt.Errorf("entity %s already exists", opt.EntityName)
			}
		}
		return 0, err
	}

	entityID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(entityID), err
}