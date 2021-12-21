package sqlite

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) WithTransaction(ctx context.Context, fn func(context.Context, metadata.DBStore) error) error {
	//TODO implement me
	panic("implement me")
}

func (db *DB) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	//TODO implement me
	panic("implement me")
}

func (db *DB) GetEntity(ctx context.Context, id int) (*types.Entity, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) GetEntityByName(ctx context.Context, name string) (*types.Entity, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) ListEntity(ctx context.Context, entityIDs *[]int) (types.EntityList, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	//TODO implement me
	panic("implement me")
}

func (db *DB) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) GetFeatureByName(ctx context.Context, name string) (*types.Feature, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) ListFeature(ctx context.Context, opt metadata.ListFeatureOpt) (types.FeatureList, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) CreateGroup(ctx context.Context, opt metadata.CreateGroupOpt) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) UpdateGroup(ctx context.Context, opt metadata.UpdateGroupOpt) error {
	//TODO implement me
	panic("implement me")
}

func (db *DB) GetGroup(ctx context.Context, id int) (*types.Group, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) GetGroupByName(ctx context.Context, name string) (*types.Group, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) ListGroup(ctx context.Context, entityID *int, groupIDs *[]int) (types.GroupList, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int, string, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	//TODO implement me
	panic("implement me")
}

func (db *DB) GetRevision(ctx context.Context, id int) (*types.Revision, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) GetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) ListRevision(ctx context.Context, groupID *int) (types.RevisionList, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) CacheListFeature(ctx context.Context, opt metadata.ListFeatureOpt) types.FeatureList {
	//TODO implement me
	panic("implement me")
}

func (db *DB) Refresh() error {
	//TODO implement me
	panic("implement me")
}
