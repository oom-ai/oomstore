package offline

import (
	"context"
	"io"

	"github.com/onestore-ai/onestore/internal/database/dbutil"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	Join(ctx context.Context, opt JoinOpt) (map[string]dbutil.RowMap, error)
	Export(ctx context.Context, opt ExportOpt) (<-chan *types.RawFeatureValueRecord, error)
	Import(ctx context.Context, opt ImportOpt) (int64, string, error)

	TypeTag(dbType string) (string, error)
	io.Closer
}
