package tikv

import (
	"context"

	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	Backend = types.BackendTiKV
)

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
	_, err := db.Client.Get(ctx, []byte(""))
	return err
}

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	panic("implement me")
}

func (db *DB) PrepareStreamTable(ctx context.Context, opt online.PrepareStreamTableOpt) error {
	return nil
}
