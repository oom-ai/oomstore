package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createRevision(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateRevisionOpt) (int, string, error) {
	var snapshotTable string
	if opt.SnapshotTable != nil {
		snapshotTable = *opt.SnapshotTable
	}

	var revisionID int
	insertQuery := "INSERT INTO feature_group_revision(group_id, revision, snapshot_table, anchored, description) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	if err := sqlxCtx.GetContext(ctx, &revisionID, insertQuery, opt.GroupID, opt.Revision, snapshotTable, opt.Anchored, opt.Description); err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return 0, "", fmt.Errorf("revision already exists: groupID=%d, revision=%d", opt.GroupID, opt.Revision)
			}
		}
		return 0, "", err
	}
	if opt.SnapshotTable == nil {
		updateQuery := "UPDATE feature_group_revision SET snapshot_table = $1 WHERE id = $2"
		snapshotTable = fmt.Sprintf("offline_%d_%d", opt.GroupID, revisionID)
		result, err := sqlxCtx.ExecContext(ctx, updateQuery, snapshotTable, revisionID)
		if err != nil {
			return 0, "", err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return 0, "", err
		}
		if rowsAffected != 1 {
			return 0, "", fmt.Errorf("failed to update revision %d: revision not found", revisionID)
		}
	}

	return revisionID, snapshotTable, nil
}
