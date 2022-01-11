package sqlite

import (
	"context"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadata/informer"
	"github.com/oom-ai/oomstore/internal/database/metadata/sqlutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/metadata"

	"github.com/jmoiron/sqlx"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
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

func Open(ctx context.Context, opt *types.SQLiteOpt) (*DB, error) {
	db, err := dbutil.OpenSQLite(opt.DBFile)
	if err != nil {
		return nil, err
	}
	informer, err := informer.New(time.Second, func() (*informer.Cache, error) {
		return sqlutil.ListMetadata(ctx, db)
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return &DB{
		DB:       db,
		Informer: informer,
	}, err
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
