package online

import (
	"context"

	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/internal/database/online/postgres"
	"github.com/onestore-ai/onestore/internal/database/online/redis"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	GetFeatureValues(ctx context.Context, dataTable, entityName, entityKey string, featureNames []string) (database.RowMap, error)
	GetFeatureValuesWithMultiEntityKeys(ctx context.Context, dataTable, entityName string, entityKeys, featureNames []string) (map[string]database.RowMap, error)
	SinkFeatureValuesStream(ctx context.Context, stream <-chan []interface{}, features []*types.Feature, revision *types.Revision) error
}

var _ Store = &postgres.DB{}
var _ Store = &redis.DB{}

type OnlineStoreOpt struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
}

func OpenPostgresDB(opt OnlineStoreOpt) (*postgres.DB, error) {
	return postgres.OpenWith(opt.Host, opt.Port, opt.User, opt.Pass, opt.Database)
}

// TODO: implement OpenRedisDB
func OpenRedisDB(opt OnlineStoreOpt) (*redis.DB, error) {
	return nil, nil
}
