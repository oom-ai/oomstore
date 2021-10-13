package database

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) GetFeature(ctx context.Context, featureName string) (*types.Feature, error) {
	var feature types.Feature
	query := `SELECT * FROM feature WHERE name = ?`
	if err := db.GetContext(ctx, &feature, query, featureName); err != nil {
		return nil, err
	}
	return &feature, nil
}
