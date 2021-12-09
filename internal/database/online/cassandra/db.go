package cassandra

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/online"
)

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	panic("implement me!")
}

func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	panic("implement me!")
}

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	panic("implement me!")
}

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	panic("implement me!")
}
