package offline

import (
	"context"
	"io"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

type Store interface {
	Join(ctx context.Context, opt JoinOpt) (*types.JoinResult, error)
	Export(ctx context.Context, opt ExportOpt) (<-chan types.ExportRecord, <-chan error)
	Import(ctx context.Context, opt ImportOpt) (int64, error)

	TypeTag(dbType string) (string, error)

	Ping(ctx context.Context) error
	io.Closer
}
