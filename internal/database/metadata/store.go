package metadata

import (
	"context"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type Store interface {
	// entity
	CreateEntity(ctx context.Context, opt types.CreateEntityOpt) error
	GetEntity(ctx context.Context, name string) (*types.Entity, error)
	ListEntity(ctx context.Context) ([]*types.Entity, error)
	UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) (int64, error)

	// feature
	CreateFeature(ctx context.Context, opt CreateFeatureOpt) error
	GetFeature(ctx context.Context, featureName string) (*types.Feature, error)
	ListFeature(ctx context.Context, opt types.ListFeatureOpt) (types.FeatureList, error)
	UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) (int64, error)

	// feature group
	CreateFeatureGroup(ctx context.Context, opt CreateFeatureGroupOpt) error
	GetFeatureGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error)
	ListFeatureGroup(ctx context.Context, entityName *string) ([]*types.FeatureGroup, error)
	UpdateFeatureGroup(ctx context.Context, opt types.UpdateFeatureGroupOpt) (int64, error)

	// revision
	CreateRevision(ctx context.Context, opt CreateRevisionOpt) (*types.Revision, error)
	ListRevision(ctx context.Context, opt ListRevisionOpt) ([]*types.Revision, error)
	GetRevision(ctx context.Context, opt GetRevisionOpt) (*types.Revision, error)
	GetLatestRevision(ctx context.Context, groupName string) (*types.Revision, error)
	UpdateRevision(ctx context.Context, opt UpdateRevisionOpt) (int64, error)
	BuildRevisionRanges(ctx context.Context, groupName string) ([]*types.RevisionRange, error)

	// transaction
	WithTransaction(ctx context.Context, fn TxFn) error

	io.Closer
}

type ExtContext interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	sqlx.ExtContext
}

type TxFn func(ctx context.Context, txStore Store) error
