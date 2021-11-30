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
	CacheGetEntity(ctx context.Context, id int) (*types.Entity, error)
	CacheGetEntityByName(ctx context.Context, name string) (*types.Entity, error)
	CacheListEntity(ctx context.Context) types.EntityList
	ListEntity(ctx context.Context) (types.EntityList, error)

	GetFeature(ctx context.Context, id int) (*types.Feature, error)
	GetFeatureByName(ctx context.Context, name string) (*types.Feature, error)
	CacheGetFeature(ctx context.Context, id int) (*types.Feature, error)
	CacheGetFeatureByName(ctx context.Context, name string) (*types.Feature, error)
	CacheListFeature(ctx context.Context, opt ListFeatureOpt) types.FeatureList

	CacheGetGroup(ctx context.Context, id int) (*types.Group, error)
	CacheGetGroupByName(ctx context.Context, name string) (*types.Group, error)
	CacheListGroup(ctx context.Context, entityID *int) types.GroupList

	GetGroup(ctx context.Context, id int) (*types.Group, error)
	GetGroupByName(ctx context.Context, name string) (*types.Group, error)
	ListGroup(ctx context.Context, entityID *int) (types.GroupList, error)

	CacheGetRevision(ctx context.Context, id int) (*types.Revision, error)
	CacheGetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error)
	CacheListRevision(ctx context.Context, groupID *int) types.RevisionList

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
	CreateGroup(ctx context.Context, opt CreateGroupOpt) (int, error)
	UpdateGroup(ctx context.Context, opt UpdateGroupOpt) error

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
