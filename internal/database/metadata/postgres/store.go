package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata/sqlutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/informer"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var _ metadata.Store = &DB{}
var _ metadata.DBStore = &Tx{}

type DB struct {
	*sqlx.DB
	*informer.Informer
}

type Tx struct {
	*sqlx.Tx
}

func Open(ctx context.Context, option *types.PostgresOpt) (*DB, error) {
	db, err := dbutil.OpenPostgresDB(option.Host, option.Port, option.User, option.Password, option.Database)
	if err != nil {
		return nil, err
	}

	// TODO: make the interval configurable
	informer, err := informer.New(time.Second, func() (*informer.Cache, error) {
		return sqlutil.ListMetadata(ctx, db)
	})
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return &DB{
		DB:       db,
		Informer: informer,
	}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	return db.PingContext(ctx)
}

func (db *DB) Close() error {
	if err := db.Informer.Close(); err != nil {
		return err
	}
	if err := db.DB.Close(); err != nil {
		return errdefs.WithStack(err)
	}
	return nil
}
