package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateEntityOpt) (int, error) {
	var entityID int
	query := "INSERT INTO entity(name, length, description) VALUES($1, $2, $3) returning id"
	err := sqlxCtx.GetContext(ctx, &entityID, query, opt.EntityName, opt.Length, opt.Description)
	if er, ok := err.(*pq.Error); ok {
		if er.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("entity %s already exists", opt.EntityName)
		}
	}
	return entityID, err
}
