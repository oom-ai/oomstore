package sqlite

import (
	"context"

	"github.com/pkg/errors"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func createGroup(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateGroupOpt) (int, error) {
	if opt.Category != types.CategoryBatch && opt.Category != types.CategoryStream {
		return 0, errdefs.InvalidAttribute(errors.Errorf("illegal category '%s', should be either 'stream' or 'batch'", opt.Category))
	}

	query := "INSERT INTO feature_group(name, entity_id, category, description) VALUES(?, ?, ?, ?)"
	res, err := sqlxCtx.ExecContext(ctx, query, opt.GroupName, opt.EntityID, opt.Category, opt.Description)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return 0, errors.Errorf("feature group %s already exists", opt.GroupName)
			}
		}
		return 0, errors.WithStack(err)
	}

	groupID, err := res.LastInsertId()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return int(groupID), nil
}
