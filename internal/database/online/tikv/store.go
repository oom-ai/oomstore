package tikv

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
)

const BackendType = types.BackendTiKV

var _ online.Store = &DB{}

type DB struct {
	*rawkv.Client
}

func Open(opt *types.TiKVOpt) (*DB, error) {
	db, err := rawkv.NewClient(context.Background(), opt.PdAddrs, config.DefaultConfig().Security)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.Client.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	_, err := db.Client.Get(ctx, []byte("ping"))
	return err
}

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	panic("implement me")
}

func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	panic("implement me")
}

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	panic("implement me")
}

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	panic("implement me")
}

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	panic("implement me")
}
