package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createRevisionTx(ctx context.Context, tx *sqlx.Tx, opt metadata.CreateRevisionOpt) (int32, string, error) {
	var dataTable string
	if opt.DataTable != nil {
		dataTable = *opt.DataTable
	}

	var revisionId int32
	insertQuery := "INSERT INTO feature_group_revision(group_id, revision, data_table, anchored, description) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	if err := tx.GetContext(ctx, &revisionId, insertQuery, opt.GroupID, opt.Revision, dataTable, opt.Anchored, opt.Description); err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return revisionId, dataTable, fmt.Errorf("revision already exists: groupId=%d, revision=%d", opt.GroupID, opt.Revision)
			}
		}
		return revisionId, dataTable, err
	}
	if opt.DataTable == nil {
		updateQuery := "UPDATE feature_group_revision SET data_table = $1 WHERE id = $2"
		dataTable = fmt.Sprintf("data_%d_%d", opt.GroupID, revisionId)
		result, err := tx.ExecContext(ctx, updateQuery, dataTable, revisionId)
		if err != nil {
			return revisionId, dataTable, err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return revisionId, dataTable, err
		}
		if rowsAffected != 1 {
			return revisionId, dataTable, fmt.Errorf("failed to update revision %d: revision not found", revisionId)
		}
	}

	return revisionId, dataTable, nil
}

// UpdateRevision = MustUpdateRevision
// If fail to update any row or update more than one row, return error
func updateRevisionTx(ctx context.Context, tx *sqlx.Tx, opt metadata.UpdateRevisionOpt) error {
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
	result, err := tx.ExecContext(ctx, tx.Rebind(query), args...)
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
