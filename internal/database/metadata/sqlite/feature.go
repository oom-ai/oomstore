package sqlite

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createFeature(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateFeatureOpt) (int, error) {
	if err := opt.ValueType.Validate(); err != nil {
		return 0, err
	}
	query := "INSERT INTO feature(name, full_name, group_id, value_type, description) VALUES (?, ?, ?, ?, ?)"
	res, err := sqlxCtx.ExecContext(ctx, sqlxCtx.Rebind(query), opt.FeatureName, opt.FullName, opt.GroupID, opt.ValueType, opt.Description)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return 0, errdefs.Errorf("feature %s already exists", opt.FeatureName)
			}
		}
		return 0, errdefs.WithStack(err)
	}

	featureID, err := res.LastInsertId()
	if err != nil {
		return 0, errdefs.WithStack(err)
	}
	return int(featureID), nil
}
