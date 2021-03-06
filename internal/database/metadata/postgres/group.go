package postgres

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func createGroup(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateGroupOpt) (int, error) {
	if opt.Category != types.CategoryBatch && opt.Category != types.CategoryStream {
		return 0, errdefs.InvalidAttribute(errdefs.Errorf("illegal category '%s', should be either 'stream' or 'batch'", opt.Category))
	}
	if opt.Category == types.CategoryStream && opt.SnapshotInterval == 0 {
		return 0, errdefs.InvalidAttribute(errdefs.Errorf("the field SnapshotInterval of the stream group %s cannot be zero", opt.GroupName))
	}

	var groupID int
	query := "insert into feature_group(name, entity_id, category, snapshot_interval, description) values($1, $2, $3, $4, $5) returning id"
	err := sqlxCtx.GetContext(ctx, &groupID, query, opt.GroupName, opt.EntityID, opt.Category, opt.SnapshotInterval, opt.Description)
	if err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return 0, errdefs.Errorf("feature group %s already exists", opt.GroupName)
			}
		}
	}
	return groupID, errdefs.WithStack(err)
}
