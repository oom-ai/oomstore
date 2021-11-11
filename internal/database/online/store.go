package online

import (
	"context"
	"io"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

type Store interface {
	Get(ctx context.Context, opt GetOpt) (dbutil.RowMap, error)
	MultiGet(ctx context.Context, opt MultiGetOpt) (map[string]dbutil.RowMap, error)
	Import(ctx context.Context, opt ImportOpt) error
	Purge(ctx context.Context, revision *typesv2.Revision) error
	io.Closer
}
