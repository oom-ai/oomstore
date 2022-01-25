package metadata

import (
	"context"
	"database/sql"
	"io"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type Store interface {
	DBStore
	CacheStore

	Ping(ctx context.Context) error
	io.Closer
}

// DBStore defines the methods that a database backend store must implement.
type DBStore interface {
	// entity
	CreateEntity(ctx context.Context, opt CreateEntityOpt) (int, error)
	UpdateEntity(ctx context.Context, opt UpdateEntityOpt) error
	GetEntity(ctx context.Context, id int) (*types.Entity, error)
	GetEntityByName(ctx context.Context, name string) (*types.Entity, error)
	ListEntity(ctx context.Context, entityIDs *[]int) (types.EntityList, error)

	// feature
	CreateFeature(ctx context.Context, opt CreateFeatureOpt) (int, error)
	UpdateFeature(ctx context.Context, opt UpdateFeatureOpt) error
	GetFeature(ctx context.Context, id int) (*types.Feature, error)
	GetFeatureByName(ctx context.Context, groupName string, featureName string) (*types.Feature, error)
	ListFeature(ctx context.Context, opt ListFeatureOpt) (types.FeatureList, error)

	// feature group
	CreateGroup(ctx context.Context, opt CreateGroupOpt) (int, error)
	UpdateGroup(ctx context.Context, opt UpdateGroupOpt) error
	GetGroup(ctx context.Context, id int) (*types.Group, error)
	GetGroupByName(ctx context.Context, name string) (*types.Group, error)
	ListGroup(ctx context.Context, entityID *int, groupIDs *[]int) (types.GroupList, error)

	// revision
	CreateRevision(ctx context.Context, opt CreateRevisionOpt) (int, string, error)
	UpdateRevision(ctx context.Context, opt UpdateRevisionOpt) error
	GetRevision(ctx context.Context, id int) (*types.Revision, error)
	GetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error)
	ListRevision(ctx context.Context, groupID *int) (types.RevisionList, error)

	// transaction
	WithTransaction(ctx context.Context, fn func(context.Context, DBStore) error) error
}

// CacheStore defines methods that a memory backend store must implement.
type CacheStore interface {
	ListCachedFeature(ctx context.Context, opt ListCachedFeatureOpt) types.FeatureList
	GetCachedGroup(ctx context.Context, id int) (*types.Group, error)

	// Refresh pulls data from database and update cache.
	Refresh() error
}

type SqlxContext interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	DriverName() string
	Rebind(string) string
}
