package mysql

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/go-sql-driver/mysql"
	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createFeature(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateFeatureOpt) (int, error) {
	if err := opt.ValueType.Validate(); err != nil {
		return 0, errdefs.WithStack(err)
	}
	query := "INSERT INTO feature(name, group_id, value_type, description) VALUES (?, ?, ?, ?)"
	res, err := sqlxCtx.ExecContext(ctx, sqlxCtx.Rebind(query), opt.FeatureName, opt.GroupID, opt.ValueType, opt.Description)
	if err != nil {
		if er, ok := err.(*mysql.MySQLError); ok {
			if er.Number == ER_DUP_ENTRY {
				return 0, errdefs.Errorf("feature %s already exists", opt.FeatureName)
			}
		}
		return 0, err
	}

	featureID, err := res.LastInsertId()
	if err != nil {
		return 0, errdefs.WithStack(err)
	}
	return int(featureID), nil
}
