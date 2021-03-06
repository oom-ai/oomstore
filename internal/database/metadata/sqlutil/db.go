package sqlutil

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata/informer"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func ListMetadata(ctx context.Context, db *sqlx.DB) (*informer.Cache, error) {
	var cache *informer.Cache
	err := dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		entities := types.EntityList{}
		if err := tx.SelectContext(ctx, &entities, `SELECT * FROM entity`); err != nil {
			return errdefs.WithStack(err)
		}

		features := types.FeatureList{}
		if err := tx.SelectContext(ctx, &features, `SELECT * FROM feature`); err != nil {
			return errdefs.WithStack(err)
		}

		groups := types.GroupList{}
		if err := tx.SelectContext(ctx, &groups, `SELECT * FROM feature_group`); err != nil {
			return errdefs.WithStack(err)
		}

		cache = informer.NewCache(entities, features, groups)
		return nil
	})
	return cache, err
}
