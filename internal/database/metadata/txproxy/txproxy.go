package txproxy

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/metadata"
)

type TxProxy struct {
	BeginTxFn func(context.Context, *sql.TxOptions) (*sqlx.Tx, error)

	CreateEntityTx func(context.Context, *sqlx.Tx, metadata.CreateEntityOpt) (int16, error)
	UpdateEntityTx func(context.Context, *sqlx.Tx, metadata.UpdateEntityOpt) error

	CreateFeatureTx func(context.Context, *sqlx.Tx, metadata.CreateFeatureOpt) (int16, error)
	UpdateFeatureTx func(context.Context, *sqlx.Tx, metadata.UpdateFeatureOpt) error

	CreateFeatureGroupTx func(context.Context, *sqlx.Tx, metadata.CreateFeatureGroupOpt) (int16, error)
	UpdateFeatureGroupTx func(context.Context, *sqlx.Tx, metadata.UpdateFeatureGroupOpt) error

	CreateRevisionTx func(context.Context, *sqlx.Tx, metadata.CreateRevisionOpt) (int32, string, error)
	UpdateRevisionTx func(context.Context, *sqlx.Tx, metadata.UpdateRevisionOpt) error
}

func (t *TxProxy) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int16, error) {
	var id int16
	err := t.WithTransaction(ctx, func(tx *sqlx.Tx) (err error) {
		id, err = t.CreateEntityTx(ctx, tx, opt)
		return err
	})
	return id, err
}

func (t *TxProxy) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	return t.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		return t.UpdateEntityTx(ctx, tx, opt)
	})
}

func (t *TxProxy) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int16, error) {
	var id int16
	err := t.WithTransaction(ctx, func(tx *sqlx.Tx) (err error) {
		id, err = t.CreateFeatureTx(ctx, tx, opt)
		return err
	})
	return id, err
}

func (t *TxProxy) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	return t.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		return t.UpdateFeatureTx(ctx, tx, opt)
	})
}

func (t *TxProxy) CreateFeatureGroup(ctx context.Context, opt metadata.CreateFeatureGroupOpt) (int16, error) {
	var id int16
	err := t.WithTransaction(ctx, func(tx *sqlx.Tx) (err error) {
		id, err = t.CreateFeatureGroupTx(ctx, tx, opt)
		return err
	})
	return id, err
}

func (t *TxProxy) UpdateFeatureGroup(ctx context.Context, opt metadata.UpdateFeatureGroupOpt) error {
	return t.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		return t.UpdateFeatureGroupTx(ctx, tx, opt)
	})
}

func (t *TxProxy) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int32, string, error) {
	var id int32
	var dataTable string
	err := t.WithTransaction(ctx, func(tx *sqlx.Tx) (err error) {
		id, dataTable, err = t.CreateRevisionTx(ctx, tx, opt)
		return err
	})
	return id, dataTable, err
}

func (t *TxProxy) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	return t.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		return t.UpdateRevisionTx(ctx, tx, opt)
	})
}

func (t *TxProxy) WithTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := t.BeginTxFn(ctx, nil)
	if err != nil {
		return nil
	}
	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()
	return fn(tx)
}
