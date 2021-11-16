package postgres

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Entity
func (tx *Tx) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) error {
	return createEntity(ctx, tx, opt)
}

func (tx *Tx) GetEntity(ctx context.Context, name string) (*types.Entity, error) {
	return getEntity(ctx, tx, name)
}

func (tx *Tx) ListEntity(ctx context.Context) ([]*types.Entity, error) {
	return listEntity(ctx, tx)
}

// Feature Group
func (tx *Tx) CreateFeatureGroup(ctx context.Context, opt metadata.CreateFeatureGroupOpt) error {
	return createFeatureGroup(ctx, tx, opt)
}

func (tx *Tx) GetFeatureGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error) {
	return getFeatureGroup(ctx, tx, groupName)
}

func (tx *Tx) ListFeatureGroup(ctx context.Context, entityName *string) ([]*types.FeatureGroup, error) {
	return listFeatureGroup(ctx, tx, entityName)
}

func (tx *Tx) UpdateFeatureGroup(ctx context.Context, opt types.UpdateFeatureGroupOpt) (int64, error) {
	return updateFeatureGroup(ctx, tx, opt)
}

// Revison
func (tx *Tx) ListRevision(ctx context.Context, opt metadata.ListRevisionOpt) ([]*types.Revision, error) {
	return listRevision(ctx, tx, opt)
}

func (tx *Tx) GetRevision(ctx context.Context, opt metadata.GetRevisionOpt) (*types.Revision, error) {
	return getRevision(ctx, tx, opt)
}

func (tx *Tx) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) (int64, error) {
	return updateRevision(ctx, tx, opt)
}

func (tx *Tx) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (*types.Revision, error) {
	return createRevision(ctx, tx, opt)
}
func (tx *Tx) GetLatestRevision(ctx context.Context, groupName string) (*types.Revision, error) {
	return getLatestRevision(ctx, tx, groupName)
}

func (tx *Tx) BuildRevisionRanges(ctx context.Context, groupName string) ([]*types.RevisionRange, error) {
	return buildRevisionRanges(ctx, tx, groupName)
}

// Feature
func (tx *Tx) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) error {
	if err := tx.validateDataType(ctx, opt.DBValueType); err != nil {
		return fmt.Errorf("err when validating value_type input, details: %s", err.Error())
	}
	return createFeature(ctx, tx, opt)
}

func (tx *Tx) GetFeature(ctx context.Context, featureName string) (*types.Feature, error) {
	return getFeature(ctx, tx, featureName)
}

func (tx *Tx) ListFeature(ctx context.Context, opt types.ListFeatureOpt) (types.FeatureList, error) {
	return listFeature(ctx, tx, opt)
}

func (tx *Tx) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) (int64, error) {
	return updateFeature(ctx, tx, opt)
}
func (tx *Tx) validateDataType(ctx context.Context, dataType string) error {
	return validateDataType(ctx, tx.Tx, dataType)
}
