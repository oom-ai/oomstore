package tikv

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/pingcap/log"
	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
	"go.uber.org/zap/zapcore"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const Backend = types.BackendTiKV

var _ online.Store = &DB{}

func init() {
	// By default, TiKV logs at INFO level. Set log level to FATAL to avoid spamming
	log.SetLevel(zapcore.FatalLevel)
}

type DB struct {
	*rawkv.Client
}

func Open(opt *types.TiKVOpt) (*DB, error) {
	db, err := rawkv.NewClient(context.Background(), opt.PdAddrs, config.DefaultConfig().Security)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.Client.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	_, err := db.Client.Get(ctx, []byte(""))
	return errdefs.WithStack(err)
}

func (db *DB) CreateTable(ctx context.Context, opt online.CreateTableOpt) error {
	return nil
}
