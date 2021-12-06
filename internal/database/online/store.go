package online

import (
	"context"
	"io"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
)

type Store interface {
	Get(ctx context.Context, opt GetOpt) (dbutil.RowMap, error)
	MultiGet(ctx context.Context, opt MultiGetOpt) (map[string]dbutil.RowMap, error)
	Import(ctx context.Context, opt ImportOpt) error
	Purge(ctx context.Context, revisionID int) error

	Ping(ctx context.Context) error
	io.Closer
}
