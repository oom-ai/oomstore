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
	CreateTable(ctx context.Context, opt CreateTableOpt) error

	// Import batch / streaming features to online store
	Import(ctx context.Context, opt ImportOpt) error

	// Push streaming feature to online store
	// Note: Make sure that the table corresponding to the stream feature already exists before executing this method
	Push(ctx context.Context, opt PushOpt) error

	Ping(ctx context.Context) error
	io.Closer
}
