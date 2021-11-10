package metadatav2

import (
	"context"
	"io"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

type Store interface {
	// entity
	CreateEntity(ctx context.Context, opt types.CreateEntityOpt) (int16, error)
	UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) (int64, error)
	GetEntity(ctx context.Context, name string) *typesv2.Entity
	ListEntity(ctx context.Context) typesv2.EntityList

	// feature
	CreateFeature(ctx context.Context, opt CreateFeatureOpt) error
	UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) (int64, error)
	GetFeature(ctx context.Context, featureName string) *typesv2.Feature
	ListFeature(ctx context.Context, opt types.ListFeatureOpt) typesv2.FeatureList

	// feature group
	CreateFeatureGroup(ctx context.Context, opt CreateFeatureGroupOpt) error
	UpdateFeatureGroup(ctx context.Context, opt types.UpdateFeatureGroupOpt) (int64, error)
	GetFeatureGroup(ctx context.Context, groupName string) *typesv2.FeatureGroup
	ListFeatureGroup(ctx context.Context, entityName *string) typesv2.FeatureGroupList

	// revision
	CreateRevision(ctx context.Context, opt CreateRevisionOpt) (*typesv2.Revision, error)
	UpdateRevision(ctx context.Context, opt UpdateRevisionOpt) (int64, error)
	ListRevision(ctx context.Context, opt ListRevisionOpt) typesv2.RevisionList
	GetRevision(ctx context.Context, opt GetRevisionOpt) (*typesv2.Revision, error)
	GetLatestRevision(ctx context.Context, groupName string) *typesv2.Revision
	BuildRevisionRanges(ctx context.Context, groupName string) []*types.RevisionRange

	io.Closer
}
