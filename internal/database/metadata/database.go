package metadata

import (
	"context"
	"io"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	// entity
	CreateEntity(ctx context.Context, opt types.CreateEntityOpt) error
	GetEntity(ctx context.Context, name string) (*types.Entity, error)
	ListEntity(ctx context.Context) ([]*types.Entity, error)
	UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) error

	// feature
	CreateFeature(ctx context.Context, opt CreateFeatureOpt) error
	GetFeature(ctx context.Context, featureName string) (*types.Feature, error)
	ListFeature(ctx context.Context, groupName *string) ([]*types.Feature, error)
	UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) error

	// rich feature
	GetRichFeature(ctx context.Context, featureName string) (*types.RichFeature, error)
	GetRichFeatures(ctx context.Context, featureNames []string) ([]*types.RichFeature, error)
	ListRichFeature(ctx context.Context, opt types.ListFeatureOpt) ([]*types.RichFeature, error)

	// feature group
	CreateFeatureGroup(ctx context.Context, opt CreateFeatureGroupOpt) error
	GetFeatureGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error)
	ListFeatureGroup(ctx context.Context, entityName *string) ([]*types.FeatureGroup, error)
	UpdateFeatureGroup(ctx context.Context, opt types.UpdateFeatureGroupOpt) error

	// revision
	ListRevision(ctx context.Context, groupName *string) ([]*types.Revision, error)
	GetRevision(ctx context.Context, groupName string, revision int64) (*types.Revision, error)
	GetRevisionsByDataTables(ctx context.Context, dataTables []string) ([]*types.Revision, error)
	GetLatestRevision(ctx context.Context, groupName string) (*types.Revision, error)
	BuildRevisionRanges(ctx context.Context, groupName string) ([]*types.RevisionRange, error)
	InsertRevision(ctx context.Context, opt InsertRevisionOpt) error

	io.Closer
}
