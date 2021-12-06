package mysql

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/metadata"
	"github.com/ethhte88/oomstore/internal/database/metadata/sqlutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func (tx *Tx) WithTransaction(ctx context.Context, fn func(context.Context, metadata.DBStore) error) error {
	return fn(ctx, tx)
}

func (tx *Tx) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int, error) {
	return createEntity(ctx, tx, opt)
}

func (tx *Tx) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	return sqlutil.UpdateEntity(ctx, tx, opt)
}

func (tx *Tx) GetEntity(ctx context.Context, id int) (*types.Entity, error) {
	return sqlutil.GetEntity(ctx, tx, id)
}

func (tx *Tx) GetEntityByName(ctx context.Context, name string) (*types.Entity, error) {
	return sqlutil.GetEntityByName(ctx, tx, name)
}

func (tx *Tx) ListEntity(ctx context.Context, entityIDs *[]int) (types.EntityList, error) {
	return sqlutil.ListEntity(ctx, tx, entityIDs)
}

func (tx *Tx) CreateGroup(ctx context.Context, opt metadata.CreateGroupOpt) (int, error) {
	return createGroup(ctx, tx, opt)
}

func (tx *Tx) UpdateGroup(ctx context.Context, opt metadata.UpdateGroupOpt) error {
	return sqlutil.UpdateGroup(ctx, tx, opt)
}

func (tx *Tx) GetGroup(ctx context.Context, id int) (*types.Group, error) {
	return sqlutil.GetGroup(ctx, tx, id)
}

func (tx *Tx) GetGroupByName(ctx context.Context, name string) (*types.Group, error) {
	return sqlutil.GetGroupByName(ctx, tx, name)
}

func (tx *Tx) ListGroup(ctx context.Context, entityID *int, groupIDs *[]int) (types.GroupList, error) {
	return sqlutil.ListGroup(ctx, tx, entityID, groupIDs)
}

func (tx *Tx) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int, error) {
	return createFeature(ctx, tx, opt)
}

func (tx *Tx) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	return sqlutil.UpdateFeature(ctx, tx, opt)
}

func (tx *Tx) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	return sqlutil.GetFeature(ctx, tx, id)
}

func (tx *Tx) ListFeature(ctx context.Context, opt metadata.ListFeatureOpt) (types.FeatureList, error) {
	return sqlutil.ListFeature(ctx, tx, opt)
}

func (tx *Tx) GetFeatureByName(ctx context.Context, name string) (*types.Feature, error) {
	return sqlutil.GetFeatureByName(ctx, tx, name)
}

func (tx *Tx) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int, string, error) {
	return createRevision(ctx, tx, opt)
}

func (tx *Tx) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	return sqlutil.UpdateRevision(ctx, tx, opt)
}

func (tx *Tx) GetRevision(ctx context.Context, id int) (*types.Revision, error) {
	return sqlutil.GetRevision(ctx, tx, id)
}

func (tx *Tx) GetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error) {
	return sqlutil.GetRevisionBy(ctx, tx, groupID, revision)
}

func (tx *Tx) ListRevision(ctx context.Context, groupID *int) (types.RevisionList, error) {
	return sqlutil.ListRevision(ctx, tx, groupID)
}
