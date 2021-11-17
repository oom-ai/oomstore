package postgres

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func (tx *Tx) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int, error) {
	return createEntity(ctx, tx, opt)
}

func (tx *Tx) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	return updateEntity(ctx, tx, opt)
}

func (tx *Tx) CreateGroup(ctx context.Context, opt metadata.CreateGroupOpt) (int, error) {
	return createGroup(ctx, tx, opt)
}

func (tx *Tx) UpdateGroup(ctx context.Context, opt metadata.UpdateGroupOpt) error {
	return updateGroup(ctx, tx, opt)
}

func (tx *Tx) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int, error) {
	return createFeature(ctx, tx, opt)
}

func (tx *Tx) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	return updateFeature(ctx, tx, opt)
}

func (tx *Tx) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int, string, error) {
	return createRevision(ctx, tx, opt)
}

func (tx *Tx) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	return updateRevision(ctx, tx, opt)
}
