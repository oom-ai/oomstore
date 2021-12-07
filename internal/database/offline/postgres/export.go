package postgres

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/internal/database/offline/sqlutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func (db *DB) Export(ctx context.Context, opt offline.ExportOpt) (<-chan types.ExportRecord, <-chan error) {
	return sqlutil.Export(ctx, db.DB, opt, types.POSTGRES)
}
