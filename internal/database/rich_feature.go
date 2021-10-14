package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) GetRichFeature(ctx context.Context, featureName string) (*types.RichFeature, error) {
	var richFeature types.RichFeature
	query := `SELECT * FROM rich_feature WHERE name = ?`
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
