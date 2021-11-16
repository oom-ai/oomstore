package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Entity
func (db *DB) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) error {
	return createEntity(ctx, db, opt)
}

func (db *DB) GetEntity(ctx context.Context, name string) (*types.Entity, error) {
	return getEntity(ctx, db, name)
}

func (db *DB) ListEntity(ctx context.Context) ([]*types.Entity, error) {
	return listEntity(ctx, db)
}

// Feature Group
func (db *DB) CreateFeatureGroup(ctx context.Context, opt metadata.CreateFeatureGroupOpt) error {
	return createFeatureGroup(ctx, db, opt)
}

func (db *DB) GetFeatureGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error) {
	return getFeatureGroup(ctx, db, groupName)
}

func (db *DB) ListFeatureGroup(ctx context.Context, entityName *string) ([]*types.FeatureGroup, error) {
	return listFeatureGroup(ctx, db, entityName)
}

func (db *DB) UpdateFeatureGroup(ctx context.Context, opt types.UpdateFeatureGroupOpt) (int64, error) {
	return updateFeatureGroup(ctx, db, opt)
}

// Revison
func (db *DB) ListRevision(ctx context.Context, opt metadata.ListRevisionOpt) ([]*types.Revision, error) {
	return listRevision(ctx, db, opt)
}
func (db *DB) GetRevision(ctx context.Context, opt metadata.GetRevisionOpt) (*types.Revision, error) {
	return getRevision(ctx, db, opt)
}
func (db *DB) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) (int64, error) {
	return updateRevision(ctx, db, opt)
}
func (db *DB) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (*types.Revision, error) {
	return createRevision(ctx, db, opt)
}
func (db *DB) GetLatestRevision(ctx context.Context, groupName string) (*types.Revision, error) {
	return getLatestRevision(ctx, db, groupName)
}

func (db *DB) BuildRevisionRanges(ctx context.Context, groupName string) ([]*types.RevisionRange, error) {
	return buildRevisionRanges(ctx, db, groupName)
}

// Feature
func (db *DB) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) error {
	if err := db.validateDataType(ctx, opt.DBValueType); err != nil {
		return fmt.Errorf("err when validating value_type input, details: %s", err.Error())
	}
	return createFeature(ctx, db, opt)
}

func (db *DB) GetFeature(ctx context.Context, featureName string) (*types.Feature, error) {
	return getFeature(ctx, db, featureName)
}

func (db *DB) ListFeature(ctx context.Context, opt types.ListFeatureOpt) (types.FeatureList, error) {
	return listFeature(ctx, db, opt)
}

func (db *DB) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) (int64, error) {
	return updateFeature(ctx, db, opt)
}

func (db *DB) validateDataType(ctx context.Context, dataType string) error {
	return dbutil.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		return validateDataType(ctx, tx, dataType)
	})
}
