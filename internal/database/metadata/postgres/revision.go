package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createRevision(ctx context.Context, exec metadata.ExecContext, opt metadata.CreateRevisionOpt) (int32, string, error) {
	var dataTable string
	if opt.DataTable != nil {
		dataTable = *opt.DataTable
	}

	var revisionId int32
	insertQuery := "INSERT INTO feature_group_revision(group_id, revision, data_table, anchored, description) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	if err := exec.GetContext(ctx, &revisionId, insertQuery, opt.GroupID, opt.Revision, dataTable, opt.Anchored, opt.Description); err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return 0, "", fmt.Errorf("revision already exists: groupId=%d, revision=%d", opt.GroupID, opt.Revision)
			}
		}
		return 0, "", err
	}
	if opt.DataTable == nil {
		updateQuery := "UPDATE feature_group_revision SET data_table = $1 WHERE id = $2"
		dataTable = fmt.Sprintf("data_%d_%d", opt.GroupID, revisionId)
		result, err := exec.ExecContext(ctx, updateQuery, dataTable, revisionId)
		if err != nil {
			return 0, "", err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return 0, "", err
		}
		if rowsAffected != 1 {
			return 0, "", fmt.Errorf("failed to update revision %d: revision not found", revisionId)
		}
	}

	return revisionId, dataTable, nil
}

// UpdateRevision = MustUpdateRevision
// If fail to update any row or update more than one row, return error
func updateRevision(ctx context.Context, exec metadata.ExecContext, opt metadata.UpdateRevisionOpt) error {
	and := make(map[string]interface{})
	if opt.NewRevision != nil {
		and["revision"] = *opt.NewRevision
	}
	if opt.NewAnchored != nil {
		and["anchored"] = *opt.NewAnchored
	}
	cond, args, err := dbutil.BuildConditions(and, nil)
	if err != nil {
		return err
	}
	if len(cond) == 0 {
		return fmt.Errorf("invliad option: nothing to update")
	}
	args = append(args, opt.RevisionID)

	query := fmt.Sprintf("UPDATE feature_group_revision SET %s WHERE id = ?", strings.Join(cond, ","))
	result, err := exec.ExecContext(ctx, exec.Rebind(query), args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("failed to update revision %d: revision not found", opt.RevisionID)
	}
	return nil
}
