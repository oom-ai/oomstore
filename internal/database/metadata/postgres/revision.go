package postgres

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createRevision(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateRevisionOpt) (int, string, error) {
	var snapshotTable, cdcTable string
	if opt.SnapshotTable != nil {
		snapshotTable = *opt.SnapshotTable
	}
	if opt.CdcTable != nil {
		cdcTable = *opt.CdcTable
	}

	var revisionID int
	insertQuery := "INSERT INTO feature_group_revision(group_id, revision, snapshot_table, cdc_table, anchored, description) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	if err := sqlxCtx.GetContext(ctx, &revisionID, insertQuery, opt.GroupID, opt.Revision, snapshotTable, cdcTable, opt.Anchored, opt.Description); err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return 0, "", errdefs.Errorf("revision already exists: groupID=%d, revision=%d", opt.GroupID, opt.Revision)
			}
		}
		return 0, "", errdefs.WithStack(err)
	}

	return revisionID, snapshotTable, nil
}
