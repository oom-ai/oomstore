package database

import (
	"context"
	"database/sql"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) ListRevision(ctx context.Context, groupName string) ([]*types.Revision, error) {
	query := "SELECT * FROM feature_group_revision WHERE group_name = ?"
	revisions := make([]*types.Revision, 0)

	if err := db.SelectContext(ctx, &revisions, query, groupName); err != nil {
		return nil, err
	}
	return revisions, nil
}

func InsertRevision(ctx context.Context, tx *sql.Tx, groupName string, revision int64, dataTable string, description string) error {
	cmd := "INSERT INTO feature_group_revision(group_name, revision, data_table, description) VALUES (?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, cmd, groupName, revision, dataTable, description)
	return err
}
