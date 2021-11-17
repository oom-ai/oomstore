package metadata

import (
	"context"
	"database/sql"
	"io"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type Store interface {
	ReadStore
	WriteStore

	io.Closer
}

type ReadStore interface {
	GetEntity(ctx context.Context, id int) (*types.Entity, error)
	GetEntityByName(ctx context.Context, name string) (*types.Entity, error)
	ListEntity(ctx context.Context) types.EntityList

	GetFeature(ctx context.Context, id int) (*types.Feature, error)
	GetFeatureByName(ctx context.Context, name string) (*types.Feature, error)
	ListFeature(ctx context.Context, opt ListFeatureOpt) types.FeatureList

	GetFeatureGroup(ctx context.Context, id int) (*types.FeatureGroup, error)
	GetFeatureGroupByName(ctx context.Context, name string) (*types.FeatureGroup, error)
	ListFeatureGroup(ctx context.Context, entityID *int) types.FeatureGroupList

	GetRevision(ctx context.Context, id int) (*types.Revision, error)
	GetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error)
	ListRevision(ctx context.Context, groupID *int) types.RevisionList

	// refresh
	Refresh() error
}

type WriteStore interface {
	// entity
	CreateEntity(ctx context.Context, opt CreateEntityOpt) (int, error)
	UpdateEntity(ctx context.Context, opt UpdateEntityOpt) error

	// feature
	CreateFeature(ctx context.Context, opt CreateFeatureOpt) (int, error)
	UpdateFeature(ctx context.Context, opt UpdateFeatureOpt) error

	// feature group
	CreateFeatureGroup(ctx context.Context, opt CreateFeatureGroupOpt) (int, error)
	UpdateFeatureGroup(ctx context.Context, opt UpdateFeatureGroupOpt) error

	// revision
	CreateRevision(ctx context.Context, opt CreateRevisionOpt) (int, string, error)
	UpdateRevision(ctx context.Context, opt UpdateRevisionOpt) error

	// transaction
	WithTransaction(ctx context.Context, fn func(context.Context, WriteStore) error) error
}

type SqlxContext interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	DriverName() string
	Rebind(string) string
}
