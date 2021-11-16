package metadata

import (
	"context"
	"io"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type Store interface {
	io.Closer
	StoreRead
	StoreWrite
}

type StoreRead interface {
	GetEntity(ctx context.Context, id int16) (*types.Entity, error)
	GetEntityByName(ctx context.Context, name string) (*types.Entity, error)
	ListEntity(ctx context.Context) types.EntityList

	GetFeature(ctx context.Context, id int16) (*types.Feature, error)
	GetFeatureByName(ctx context.Context, name string) (*types.Feature, error)
	ListFeature(ctx context.Context, opt ListFeatureOpt) types.FeatureList

	GetFeatureGroup(ctx context.Context, id int16) (*types.FeatureGroup, error)
	GetFeatureGroupByName(ctx context.Context, name string) (*types.FeatureGroup, error)
	ListFeatureGroup(ctx context.Context, entityID *int16) types.FeatureGroupList

	GetRevision(ctx context.Context, id int32) (*types.Revision, error)
	GetRevisionBy(ctx context.Context, groupID int16, revision int64) (*types.Revision, error)
	ListRevision(ctx context.Context, opt ListRevisionOpt) types.RevisionList

	Refresh() error
}

type StoreWrite interface {
	// entity
	CreateEntity(ctx context.Context, opt CreateEntityOpt) (int16, error)
	UpdateEntity(ctx context.Context, opt UpdateEntityOpt) error

	// feature
	CreateFeature(ctx context.Context, opt CreateFeatureOpt) (int16, error)
	UpdateFeature(ctx context.Context, opt UpdateFeatureOpt) error

	// feature group
	CreateFeatureGroup(ctx context.Context, opt CreateFeatureGroupOpt) (int16, error)
	UpdateFeatureGroup(ctx context.Context, opt UpdateFeatureGroupOpt) error

	// revision
	CreateRevision(ctx context.Context, opt CreateRevisionOpt) (int32, string, error)
	UpdateRevision(ctx context.Context, opt UpdateRevisionOpt) error
}
