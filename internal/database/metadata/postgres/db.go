package postgres

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func (db *DB) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int, error) {
	return createEntity(ctx, db, opt)
}

func (db *DB) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	return updateEntity(ctx, db, opt)
}

func (db *DB) CreateFeatureGroup(ctx context.Context, opt metadata.CreateFeatureGroupOpt) (int, error) {
	return createFeatureGroup(ctx, db, opt)
}

func (db *DB) UpdateFeatureGroup(ctx context.Context, opt metadata.UpdateFeatureGroupOpt) error {
	return updateFeatureGroup(ctx, db, opt)
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
	err = db.WithTransaction(ctx, func(c context.Context, s metadata.StoreWrite) error {
		revisionID, dataTable, err = createRevision(ctx, db, opt)
		return err
	})
	return revisionID, dataTable, err
}

func (db *DB) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	return updateRevision(ctx, db, opt)
}
