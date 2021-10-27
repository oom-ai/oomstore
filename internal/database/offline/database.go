package offline

import (
	"context"
	"io"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type Store interface {
	Join(ctx context.Context, opt JoinOpt) (map[string]dbutil.RowMap, error)
	Export(ctx context.Context, opt ExportOpt) (<-chan *types.RawFeatureValueRecord, error)
	Import(ctx context.Context, opt ImportOpt) (int64, string, error)

	TypeTag(dbType string) (string, error)
	io.Closer
}
