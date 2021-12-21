package sqlite

import (
	"context"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createRevision(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateRevisionOpt) (int, string, error) {
	var dataTable string
	if opt.DataTable != nil {
		dataTable = *opt.DataTable
	}

	insertQuery := "INSERT INTO feature_group_revision(group_id, revision, data_table, anchored, description) VALUES (?, ?, ?, ?, ?)"
	res, err := sqlxCtx.ExecContext(ctx, sqlxCtx.Rebind(insertQuery), opt.GroupID, opt.Revision, dataTable, opt.Anchored, opt.Description)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, "", fmt.Errorf("revision already exists: groupID=%d, revision=%d", opt.GroupID, opt.Revision)
			}
		}
		return 0, "", err
	}
	revisionID, err := res.LastInsertId()
	if err != nil {
		return 0, "", err
	}

	if opt.DataTable == nil {
		updateQuery := "UPDATE feature_group_revision SET data_table = ? WHERE id = ?"
		dataTable = fmt.Sprintf("offline_%d_%d", opt.GroupID, revisionID)
		result, err := sqlxCtx.ExecContext(ctx, sqlxCtx.Rebind(updateQuery), dataTable, revisionID)
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
	return int(revisionID), dataTable, nil
}
