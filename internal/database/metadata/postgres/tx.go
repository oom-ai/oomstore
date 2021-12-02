package postgres

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

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
	return updateGroup(ctx, tx, opt)
}

func (tx *Tx) GetGroup(ctx context.Context, id int) (*types.Group, error) {
	return getGroup(ctx, tx, id)
}

func (tx *Tx) GetGroupByName(ctx context.Context, name string) (*types.Group, error) {
	return getGroupByName(ctx, tx, name)
}

func (tx *Tx) ListGroup(ctx context.Context, entityID *int, groupIDs *[]int) (types.GroupList, error) {
	return listGroup(ctx, tx, entityID, groupIDs)
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

func (tx *Tx) ListFeature(ctx context.Context, opt metadata.ListFeatureOpt) (types.FeatureList, error) {
	return listFeature(ctx, tx, opt)
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

func (tx *Tx) GetRevision(ctx context.Context, id int) (*types.Revision, error) {
	return getRevision(ctx, tx, id)
}

func (tx *Tx) GetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error) {
	return getRevisionBy(ctx, tx, groupID, revision)
}

func (tx *Tx) ListRevision(ctx context.Context, groupID *int) (types.RevisionList, error) {
	return listRevision(ctx, tx, groupID)
}
