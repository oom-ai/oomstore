package metadata

import (
	"context"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type Store interface {
	// entity
	CreateEntity(ctx context.Context, opt CreateEntityOpt) (int16, error)
	UpdateEntity(ctx context.Context, opt UpdateEntityOpt) error
	GetEntity(ctx context.Context, id int16) (*types.Entity, error)
	GetEntityByName(ctx context.Context, name string) (*types.Entity, error)
	ListEntity(ctx context.Context) types.EntityList

	// feature
	CreateFeature(ctx context.Context, opt CreateFeatureOpt) (int16, error)
	UpdateFeature(ctx context.Context, opt UpdateFeatureOpt) error
	GetFeature(ctx context.Context, id int16) (*types.Feature, error)
	GetFeatureByName(ctx context.Context, name string) (*types.Feature, error)
	ListFeature(ctx context.Context, opt ListFeatureOpt) types.FeatureList

	// feature group
	CreateFeatureGroup(ctx context.Context, opt CreateFeatureGroupOpt) (int16, error)
	UpdateFeatureGroup(ctx context.Context, opt UpdateFeatureGroupOpt) error
	GetFeatureGroup(ctx context.Context, id int16) (*types.FeatureGroup, error)
	GetFeatureGroupByName(ctx context.Context, name string) (*types.FeatureGroup, error)
	ListFeatureGroup(ctx context.Context, entityID *int16) types.FeatureGroupList

	// revision
	CreateRevision(ctx context.Context, opt CreateRevisionOpt) (int32, error)
	UpdateRevision(ctx context.Context, opt UpdateRevisionOpt) error
	GetRevision(ctx context.Context, id int32) (*types.Revision, error)
	GetRevisionBy(ctx context.Context, groupID int16, revision int64) (*types.Revision, error)
	ListRevision(ctx context.Context, opt ListRevisionOpt) types.RevisionList

	// transaction
	WithTransaction(context.Context, func(context.Context, Store) error) error
	// refresh
	Refresh() error
	io.Closer
}

type ExtContext interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	sqlx.ExtContext
}
