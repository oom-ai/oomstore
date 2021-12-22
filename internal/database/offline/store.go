package offline

import (
	"context"
	"io"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type Store interface {
	Join(ctx context.Context, opt JoinOpt) (*types.JoinResult, error)
	Export(ctx context.Context, opt ExportOpt) (<-chan types.ExportRecord, <-chan error)
	Import(ctx context.Context, opt ImportOpt) (int64, error)

	TableSchema(ctx context.Context, tableName string) (*types.DataTableSchema, error)

	Ping(ctx context.Context) error
	io.Closer
}
