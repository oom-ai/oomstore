package metadata

import (
	"context"
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *PostgresDB) ListRevision(ctx context.Context, groupName *string) ([]*types.Revision, error) {
	query := "SELECT * FROM feature_group_revision"
	var cond []interface{}
	if groupName != nil {
		query += " WHERE group_name = $1"
		cond = append(cond, *groupName)
	}
	revisions := make([]*types.Revision, 0)

	if err := db.SelectContext(ctx, &revisions, query, cond...); err != nil {
		return nil, err
	}
	return revisions, nil
}

func InsertRevision(ctx context.Context, tx *sqlx.Tx, groupName string, revision int64, dataTable string, description string) error {
	cmd := "INSERT INTO feature_group_revision(group_name, revision, data_table, description) VALUES ($1, $2, $3, $4)"
	_, err := tx.ExecContext(ctx, cmd, groupName, revision, dataTable, description)
	return err
}

func (db *PostgresDB) BuildRevisionRanges(ctx context.Context, groupName string) ([]*types.RevisionRange, error) {
	query := fmt.Sprintf(`
		SELECT
			revision AS min_revision,
			LEAD(revision, 1, %d) OVER (ORDER BY revision) AS max_revision,
			data_table
		FROM feature_group_revision
		WHERE group_name = $1
	`, math.MaxInt64)

	var ranges []*types.RevisionRange
	if err := db.SelectContext(ctx, &ranges, query, groupName); err != nil {
		return nil, err
	}
	return ranges, nil
}
