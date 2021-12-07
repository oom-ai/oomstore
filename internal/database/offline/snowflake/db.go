package snowflake

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func (db *DB) Import(ctx context.Context, opt offline.ImportOpt) (int64, error) {
	panic("implement me")
}

func (db *DB) Export(ctx context.Context, opt offline.ExportOpt) (<-chan types.ExportRecord, <-chan error) {
	panic("implement me")
}

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	panic("implement me")
}

func (db *DB) TypeTag(dbType string) (string, error) {
	panic("implement me")
}
