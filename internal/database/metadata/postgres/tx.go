package postgres

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (tx *Tx) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int, error) {
	return createEntity(ctx, tx, opt)
}

func (tx *Tx) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	return updateEntity(ctx, tx, opt)
}

func (tx *Tx) GetEntity(ctx context.Context, id int) (*types.Entity, error) {
	return getEntity(ctx, tx, id)
}

func (tx *Tx) GetEntityByName(ctx context.Context, name string) (*types.Entity, error) {
	return getEntityByName(ctx, tx, name)
}

func (tx *Tx) ListEntity(ctx context.Context, entityIDs *[]int) (types.EntityList, error) {
	return listEntity(ctx, tx, entityIDs)
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

func (tx *Tx) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	return getFeature(ctx, tx, id)
}

func (tx *Tx) GetFeatureByName(ctx context.Context, name string) (*types.Feature, error) {
	return getFeatureByName(ctx, tx, name)
}

func (tx *Tx) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int, string, error) {
	return createRevision(ctx, tx, opt)
}

func (tx *Tx) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	return updateRevision(ctx, tx, opt)
}
