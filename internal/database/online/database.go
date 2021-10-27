package online

import (
	"context"
	"fmt"
	"io"

	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/internal/database/online/postgres"
	"github.com/onestore-ai/onestore/internal/database/online/redis"
	dbtypes "github.com/onestore-ai/onestore/internal/database/types"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	Get(ctx context.Context, opt types.GetFeatureValuesOpt) (database.RowMap, error)
	MultiGet(ctx context.Context, opt dbtypes.MultiGetOnlineFeatureValuesOpt) (map[string]database.RowMap, error)
	Import(ctx context.Context, stream <-chan *types.RawFeatureValueRecord, features []*types.Feature, revision *types.Revision, entity *types.Entity) error
	PurgeRevision(ctx context.Context, revision *types.Revision) error
	io.Closer
}

var _ Store = &postgres.DB{}
var _ Store = &redis.DB{}

func Open(opt types.OnlineStoreOpt) (Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return postgres.Open(opt.PostgresDbOpt)
	case types.REDIS:
		return redis.Open(opt.RedisDbOpt), nil
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}
