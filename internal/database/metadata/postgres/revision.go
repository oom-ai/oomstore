package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/sqlutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func createRevision(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateRevisionOpt) (int, string, error) {
	var dataTable string
	if opt.DataTable != nil {
		dataTable = *opt.DataTable
	}

	var revisionID int
	insertQuery := "INSERT INTO feature_group_revision(group_id, revision, data_table, anchored, description) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	if err := sqlxCtx.GetContext(ctx, &revisionID, insertQuery, opt.GroupID, opt.Revision, dataTable, opt.Anchored, opt.Description); err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return 0, "", fmt.Errorf("revision already exists: groupID=%d, revision=%d", opt.GroupID, opt.Revision)
			}
		}
		return 0, "", err
	}
	if opt.DataTable == nil {
		updateQuery := "UPDATE feature_group_revision SET data_table = $1 WHERE id = $2"
		dataTable = fmt.Sprintf("offline_%d_%d", opt.GroupID, revisionID)
		result, err := sqlxCtx.ExecContext(ctx, updateQuery, dataTable, revisionID)
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

	return revisionID, dataTable, nil
}

// UpdateRevision = MustUpdateRevision
// If fail to update any row or update more than one row, return error
func updateRevision(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.UpdateRevisionOpt) error {
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
	result, err := sqlxCtx.ExecContext(ctx, sqlxCtx.Rebind(query), args...)
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

func getRevision(ctx context.Context, sqlxCtx metadata.SqlxContext, id int) (*types.Revision, error) {
	var revision types.Revision
	query := `SELECT * FROM "feature_group_revision" WHERE id = $1`
	if err := sqlxCtx.GetContext(ctx, &revision, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(fmt.Errorf("revision %d not found", id))
		}
		return nil, err
	}

	group, err := sqlutil.GetGroup(ctx, sqlxCtx, revision.GroupID)
	if err != nil {
		return nil, err
	}
	revision.Group = group
	return &revision, nil
}

func getRevisionBy(ctx context.Context, sqlxCtx metadata.SqlxContext, groupID int, revision int64) (*types.Revision, error) {
	var r types.Revision
	query := `SELECT * FROM "feature_group_revision" WHERE "group_id" = $1 AND "revision" = $2`
	if err := sqlxCtx.GetContext(ctx, &r, query, groupID, revision); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(fmt.Errorf("revision not found by group %d and revision %d", groupID, revision))
		}
		return nil, err
	}

	group, err := sqlutil.GetGroup(ctx, sqlxCtx, r.GroupID)
	if err != nil {
		return nil, err
	}
	r.Group = group
	return &r, nil
}
func listRevision(ctx context.Context, sqlxCtx metadata.SqlxContext, groupID *int) (types.RevisionList, error) {
	query := `SELECT * FROM "feature_group_revision"`
	args := make([]interface{}, 0)
	if groupID != nil {
		query = fmt.Sprintf(`%s WHERE group_id = $1`, query)
		args = append(args, *groupID)
	}
	var revisions types.RevisionList
	if err := sqlxCtx.SelectContext(ctx, &revisions, query, args...); err != nil {
		return nil, err
	}

	if err := enrichRevisions(ctx, sqlxCtx, revisions); err != nil {
		return nil, err
	}
	return revisions, nil
}

func enrichRevisions(ctx context.Context, sqlxCtx metadata.SqlxContext, revisions types.RevisionList) error {
	groupIDs := revisions.GroupIDs()
	groups, err := sqlutil.ListGroup(ctx, sqlxCtx, nil, &groupIDs)
	if err != nil {
		return err
	}
	for _, revision := range revisions {
		group := groups.Find(func(g *types.Group) bool {
			return revision.GroupID == g.ID
		})
		if group == nil {
			return fmt.Errorf("cannot find group %d for revision %d", revision.GroupID, revision.ID)
		}
		revision.Group = group
	}
	return nil
}
