package database

import (
	"context"

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
