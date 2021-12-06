package sqlutil

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/metadata/informer"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
)

func ListMetaData(ctx context.Context, db *sqlx.DB) (*informer.Cache, error) {
	var cache *informer.Cache
	err := dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		entities := types.EntityList{}
		if err := tx.SelectContext(ctx, &entities, `SELECT * FROM entity`); err != nil {
			return err
		}

		features := types.FeatureList{}
		if err := tx.SelectContext(ctx, &features, `SELECT * FROM feature`); err != nil {
			return err
		}

		groups := types.GroupList{}
		if err := tx.SelectContext(ctx, &groups, `SELECT * FROM feature_group`); err != nil {
			return err
		}

		cache = informer.NewCache(entities, features, groups)
		return nil
	})
	return cache, err
}
