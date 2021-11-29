package postgres

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int, error) {
	return createEntity(ctx, db, opt)
}

func (db *DB) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	return updateEntity(ctx, db, opt)
}

func (db *DB) CreateGroup(ctx context.Context, opt metadata.CreateGroupOpt) (int, error) {
	return createGroup(ctx, db, opt)
}

func (db *DB) UpdateGroup(ctx context.Context, opt metadata.UpdateGroupOpt) error {
	return updateGroup(ctx, db, opt)
}

func (db *DB) GetGroup(ctx context.Context, id int) (*types.Group, error) {
	return getGroup(ctx, db, id)
}

func (db *DB) GetGroupByName(ctx context.Context, name string) (*types.Group, error) {
	return getGroupByName(ctx, db, name)
}

func (db *DB) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int, error) {
	return createFeature(ctx, db, opt)
}

func (db *DB) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	return updateFeature(ctx, db, opt)
}

func (db *DB) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int, string, error) {
	var (
		revisionID int
		dataTable  string
		err        error
	)
	err = db.WithTransaction(ctx, func(c context.Context, s metadata.WriteStore) error {
		revisionID, dataTable, err = createRevision(ctx, db, opt)
		return err
	})
	return revisionID, dataTable, err
}

func (db *DB) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	return updateRevision(ctx, db, opt)
}
