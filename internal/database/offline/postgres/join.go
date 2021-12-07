package postgres

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/internal/database/offline/sqlutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	return sqlutil.Join(ctx, db.DB, opt, types.POSTGRES)
}
