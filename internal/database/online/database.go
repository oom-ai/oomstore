package online

import (
	"context"
	"io"

	"github.com/onestore-ai/onestore/internal/database/dbutil"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	Get(ctx context.Context, opt GetOpt) (dbutil.RowMap, error)
	MultiGet(ctx context.Context, opt MultiGetOpt) (map[string]dbutil.RowMap, error)
	Import(ctx context.Context, opt ImportOpt) error
	Purge(ctx context.Context, revision *types.Revision) error
	io.Closer
}
