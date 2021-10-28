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

func (db *DB) GetRichFeatures(ctx context.Context, featureNames []string) (types.RichFeatureList, error) {
	query := "SELECT * FROM rich_feature WHERE name IN (?)"
	sql, args, err := sqlx.In(query, featureNames)
	if err != nil {
		return nil, err
	}

	features := types.RichFeatureList{}
	err = db.SelectContext(ctx, &features, db.Rebind(sql), args...)
	if err != nil {
		return nil, err
	}
	return features, nil
}

func (db *DB) ListRichFeature(ctx context.Context, opt types.ListFeatureOpt) (types.RichFeatureList, error) {
	query := "SELECT * FROM rich_feature"
	cond, args := buildListFeatureCond(opt)
	if len(cond) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, cond)
	}

	features := types.RichFeatureList{}
	if err := db.SelectContext(ctx, &features, query, args...); err != nil {
		return nil, err
	}
	return features, nil
}
