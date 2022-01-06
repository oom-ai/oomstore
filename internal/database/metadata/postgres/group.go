package postgres

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func createGroup(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateGroupOpt) (int, error) {
	if opt.Category != types.CategoryBatch && opt.Category != types.CategoryStream {
		return 0, errdefs.InvalidAttribute(errors.Errorf("illegal category '%s', should be either 'stream' or 'batch'", opt.Category))
	}
	var groupID int
	query := "insert into feature_group(name, entity_id, category, description) values($1, $2, $3, $4) returning id"
	err := sqlxCtx.GetContext(ctx, &groupID, query, opt.GroupName, opt.EntityID, opt.Category, opt.Description)
	if err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return 0, errors.Errorf("feature group %s already exists", opt.GroupName)
			}
		}
	}
	return groupID, errors.WithStack(err)
}
