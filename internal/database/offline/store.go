package offline

import (
	"context"
	"io"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type Store interface {
	Join(ctx context.Context, opt JoinOpt) (*types.JoinResult, error)
	Export(ctx context.Context, opt ExportOpt) (*types.ExportResult, error)
	Import(ctx context.Context, opt ImportOpt) (int64, error)
	Push(ctx context.Context, opt PushOpt) error

	CreateTable(ctx context.Context, opt CreateTableOpt) error
	TableSchema(ctx context.Context, tableName string) (*types.DataTableSchema, error)
	Snapshot(ctx context.Context, opt SnapshotOpt) error

	Ping(ctx context.Context) error
	io.Closer
}
