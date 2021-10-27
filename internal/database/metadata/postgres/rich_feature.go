package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) GetRichFeature(ctx context.Context, featureName string) (*types.RichFeature, error) {
	var richFeature types.RichFeature
	query := "SELECT * FROM rich_feature WHERE name = $1"
	if err := db.GetContext(ctx, &richFeature, query, featureName); err != nil {
		return nil, err
	}
	return &richFeature, nil
}

func (db *DB) GetRichFeatures(ctx context.Context, featureNames []string) ([]*types.RichFeature, error) {
	query := "SELECT * FROM rich_feature WHERE name IN (?)"
	sql, args, err := sqlx.In(query, featureNames)
	if err != nil {
		return nil, err
	}

	richFeatures := make([]*types.RichFeature, 0)
	err = db.SelectContext(ctx, &richFeatures, db.Rebind(sql), args...)
	if err != nil {
		return nil, err
	}
	return richFeatures, nil
}

func (db *DB) ListRichFeature(ctx context.Context, opt types.ListFeatureOpt) ([]*types.RichFeature, error) {
	query := "SELECT * FROM rich_feature"
	cond, args := buildListFeatureCond(opt)
	if len(cond) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, cond)
	}

	features := make([]*types.RichFeature, 0)
	if err := db.SelectContext(ctx, &features, query, args...); err != nil {
		return nil, err
	}
	return features, nil
}
