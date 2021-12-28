package online

import (
	"context"
	"io"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
)

type Store interface {
	Get(ctx context.Context, opt GetOpt) (dbutil.RowMap, error)
	MultiGet(ctx context.Context, opt MultiGetOpt) (map[string]dbutil.RowMap, error)
	Purge(ctx context.Context, revisionID int) error

	// Batch import batch feature to online store
	Import(ctx context.Context, opt ImportOpt) error

	// Push streaming feature to online store
	Push(ctx context.Context, opt PushOpt) error

	Ping(ctx context.Context) error
	io.Closer
}
