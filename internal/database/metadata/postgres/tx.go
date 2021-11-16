package postgres

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func (tx *Tx) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int16, error) {
	return createEntity(ctx, tx, opt)
}

func (tx *Tx) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	return updateEntity(ctx, tx, opt)
}

func (tx *Tx) CreateFeatureGroup(ctx context.Context, opt metadata.CreateFeatureGroupOpt) (int16, error) {
	return createFeatureGroup(ctx, tx, opt)
}

func (tx *Tx) UpdateFeatureGroup(ctx context.Context, opt metadata.UpdateFeatureGroupOpt) error {
	return updateFeatureGroup(ctx, tx, opt)
}

func (tx *Tx) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int16, error) {
	return createFeature(ctx, tx, opt)
}

func (tx *Tx) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	return updateFeature(ctx, tx, opt)
}

func (tx *Tx) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int32, string, error) {
	return createRevision(ctx, tx, opt)
}

func (tx *Tx) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	return updateRevision(ctx, tx, opt)
}
