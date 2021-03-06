package postgres

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/errdefs"
)

func createFeature(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateFeatureOpt) (int, error) {
	if err := opt.ValueType.Validate(); err != nil {
		return 0, err
	}
	var featureID int
	query := "INSERT INTO feature(name, group_id, value_type, description) VALUES ($1, $2, $3, $4) RETURNING id"
	err := sqlxCtx.GetContext(ctx, &featureID, query, opt.FeatureName, opt.GroupID, opt.ValueType, opt.Description)
	if err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return 0, errdefs.Errorf("feature %s already exists", opt.FeatureName)
			}
		}
	}
	return featureID, errdefs.WithStack(err)
}
