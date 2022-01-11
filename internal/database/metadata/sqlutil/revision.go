package sqlutil

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// UpdateRevision = MustUpdateRevision
// If fail to update any row or update more than one row, return error
func UpdateRevision(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.UpdateRevisionOpt) error {
	and := make(map[string]interface{})
	if opt.NewRevision != nil {
		and["revision"] = *opt.NewRevision
	}
	if opt.NewAnchored != nil {
		and["anchored"] = *opt.NewAnchored
	}
	if opt.NewSnapshotTable != nil {
		and["snapshot_table"] = *opt.NewSnapshotTable
	}
	cond, args, err := dbutil.BuildConditions(and, nil)
	if err != nil {
		return err
	}
	if len(cond) == 0 {
		return errdefs.Errorf("invliad option: nothing to update")
	}
	args = append(args, opt.RevisionID)

	query := fmt.Sprintf("UPDATE feature_group_revision SET %s WHERE id = ?", strings.Join(cond, ","))
	result, err := sqlxCtx.ExecContext(ctx, sqlxCtx.Rebind(query), args...)
	if err != nil {
		return errdefs.WithStack(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errdefs.WithStack(err)
	}
	if rowsAffected != 1 {
		return errdefs.Errorf("failed to update revision %d: revision not found", opt.RevisionID)
	}
	return nil
}

func GetRevision(ctx context.Context, sqlxCtx metadata.SqlxContext, id int) (*types.Revision, error) {
	var revision types.Revision
	query := `SELECT * FROM feature_group_revision WHERE id = ?`
	if err := sqlxCtx.GetContext(ctx, &revision, sqlxCtx.Rebind(query), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(errdefs.Errorf("revision %d not found", id))
		}
		return nil, errdefs.WithStack(err)
	}

	group, err := GetGroup(ctx, sqlxCtx, revision.GroupID)
	if err != nil {
		return nil, err
	}
	revision.Group = group
	return &revision, nil
}

func GetRevisionBy(ctx context.Context, sqlxCtx metadata.SqlxContext, groupID int, revision int64) (*types.Revision, error) {
	var r types.Revision
	query := `SELECT * FROM feature_group_revision WHERE group_id = ? AND revision = ?`
	if err := sqlxCtx.GetContext(ctx, &r, sqlxCtx.Rebind(query), groupID, revision); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(errdefs.Errorf("revision not found by group %d and revision %d", groupID, revision))
		}
		return nil, errdefs.WithStack(err)
	}

	group, err := GetGroup(ctx, sqlxCtx, r.GroupID)
	if err != nil {
		return nil, err
	}
	r.Group = group
	return &r, nil
}

func ListRevision(ctx context.Context, sqlxCtx metadata.SqlxContext, groupID *int) (types.RevisionList, error) {
	query := `SELECT * FROM feature_group_revision`
	args := make([]interface{}, 0)
	if groupID != nil {
		query = fmt.Sprintf(`%s WHERE group_id = ?`, query)
		args = append(args, *groupID)
	}
	query = fmt.Sprintf("%s ORDER BY id ASC", query)
	var revisions types.RevisionList
	if err := sqlxCtx.SelectContext(ctx, &revisions, sqlxCtx.Rebind(query), args...); err != nil {
		return nil, errdefs.WithStack(err)
	}

	if err := enrichRevisions(ctx, sqlxCtx, revisions); err != nil {
		return nil, err
	}
	return revisions, nil
}

func enrichRevisions(ctx context.Context, sqlxCtx metadata.SqlxContext, revisions types.RevisionList) error {
	groupIDs := revisions.GroupIDs()
	groups, err := ListGroup(ctx, sqlxCtx, nil, &groupIDs)
	if err != nil {
		return err
	}
	for _, revision := range revisions {
		group := groups.Find(func(g *types.Group) bool {
			return revision.GroupID == g.ID
		})
		if group == nil {
			return errdefs.Errorf("cannot find group %d for revision %d", revision.GroupID, revision.ID)
		}
		revision.Group = group
	}
	return nil
}
