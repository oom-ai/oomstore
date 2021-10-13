package database

import (
	"context"
)

func (db *DB) CreateEntity(ctx context.Context, entityName, description string) error {
	_, err := db.ExecContext(ctx,
		"insert into entity(name, description) values(?, ?)",
		entityName, description)
	return err
}
