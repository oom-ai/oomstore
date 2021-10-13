package database

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) CreateGroup(ctx context.Context, opt types.CreateGroupOpt) error {
	_, err := db.ExecContext(ctx,
		"insert into feature_group(name, entity_name, category, description) values(?, ?, ?, ?, ?, ?)",
		opt.Name, opt.EntityName, opt.Category, opt.Description)
	return err
}
