package mysql

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	ER_DUP_ENTRY = 1062
)

func (db *DB) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int, error) {
	return createEntity(ctx, db, opt)
}

func (db *DB) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	return sqlutil.UpdateEntity(ctx, db, opt)
}

func (db *DB) GetEntity(ctx context.Context, id int) (*types.Entity, error) {
	return sqlutil.GetEntity(ctx, db, id)
}

func (db *DB) GetEntityByName(ctx context.Context, name string) (*types.Entity, error) {
	return sqlutil.GetEntityByName(ctx, db, name)
}

func (db *DB) ListEntity(ctx context.Context, entityIDs *[]int) (types.EntityList, error) {
	return sqlutil.ListEntity(ctx, db, entityIDs)
}

func (db *DB) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int, error) {
	panic("implement me")
}

func (db *DB) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	panic("implement me")
}

func (db *DB) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	panic("implement me")
}

func (db *DB) GetFeatureByName(ctx context.Context, name string) (*types.Feature, error) {
	panic("implement me")
}

func (db *DB) ListFeature(ctx context.Context, opt metadata.ListFeatureOpt) (types.FeatureList, error) {
	panic("implement me")
}

func (db *DB) CreateGroup(ctx context.Context, opt metadata.CreateGroupOpt) (int, error) {
	panic("implement me")
}

func (db *DB) UpdateGroup(ctx context.Context, opt metadata.UpdateGroupOpt) error {
	panic("implement me")
}

func (db *DB) GetGroup(ctx context.Context, id int) (*types.Group, error) {
	panic("implement me")
}

func (db *DB) GetGroupByName(ctx context.Context, name string) (*types.Group, error) {
	panic("implement me")
}

func (db *DB) ListGroup(ctx context.Context, entityID *int, groupIDs *[]int) (types.GroupList, error) {
	panic("implement me")
}

func (db *DB) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int, string, error) {
	panic("implement me")
}

func (db *DB) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	panic("implement me")
}

func (db *DB) GetRevision(ctx context.Context, id int) (*types.Revision, error) {
	panic("implement me")
}

func (db *DB) GetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error) {
	panic("implement me")
}

func (db *DB) ListRevision(ctx context.Context, groupID *int) (types.RevisionList, error) {
	panic("implement me")
}

func (db *DB) WithTransaction(ctx context.Context, fn func(context.Context, metadata.DBStore) error) error {
	panic("implement me")
}
