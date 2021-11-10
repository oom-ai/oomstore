package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadatav2"
)

func (db *DB) CreateRevision(ctx context.Context, opt metadatav2.CreateRevisionOpt) (int32, error) {
	query := "INSERT INTO feature_group_revision(group_id, revision, data_table, anchored, description) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var revisionId int32
	if err := db.GetContext(ctx, &revisionId, query, opt.GroupId, opt.Revision, opt.DataTable, opt.Anchored, opt.Description); err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return 0, fmt.Errorf("revision %v already exist", opt.Revision)
			}
		}
		return 0, err
	}
	return revisionId, nil
}

// UpdateRevision = MustUpdateRevision
// If fail to update any row or update more than one row, return error
func (db *DB) UpdateRevision(ctx context.Context, opt metadatav2.UpdateRevisionOpt) error {
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
	result, err := db.ExecContext(ctx, db.Rebind(query), args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("failed to update revision %d", opt.RevisionID)
	}
	return nil
}
