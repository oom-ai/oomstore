package metadata

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Tx struct {
	tx *sqlx.Tx
	tp *TxProxy
}

func (tx *Tx) CreateEntity(ctx context.Context, opt CreateEntityOpt) (int16, error) {
	return tx.tp.CreateEntityTx(ctx, tx.tx, opt)
}
func (tx *Tx) UpdateEntity(ctx context.Context, opt UpdateEntityOpt) error {
	return tx.tp.UpdateEntityTx(ctx, tx.tx, opt)
}
func (tx *Tx) CreateFeature(ctx context.Context, opt CreateFeatureOpt) (int16, error) {
	return tx.tp.CreateFeatureTx(ctx, tx.tx, opt)
}
func (tx *Tx) UpdateFeature(ctx context.Context, opt UpdateFeatureOpt) error {
	return tx.tp.UpdateFeatureTx(ctx, tx.tx, opt)
}
func (tx *Tx) CreateFeatureGroup(ctx context.Context, opt CreateFeatureGroupOpt) (int16, error) {
	return tx.tp.CreateFeatureGroupTx(ctx, tx.tx, opt)
}
func (tx *Tx) UpdateFeatureGroup(ctx context.Context, opt UpdateFeatureGroupOpt) error {
	return tx.tp.UpdateFeatureGroupTx(ctx, tx.tx, opt)
}
func (tx *Tx) CreateRevision(ctx context.Context, opt CreateRevisionOpt) (int32, string, error) {
	return tx.tp.CreateRevisionTx(ctx, tx.tx, opt)
}
func (tx *Tx) UpdateRevision(ctx context.Context, opt UpdateRevisionOpt) error {
	return tx.tp.UpdateRevisionTx(ctx, tx.tx, opt)
}

type TxProxy struct {
	BeginTxFn func(context.Context, *sql.TxOptions) (*sqlx.Tx, error)

	CreateEntityTx func(context.Context, *sqlx.Tx, CreateEntityOpt) (int16, error)
	UpdateEntityTx func(context.Context, *sqlx.Tx, UpdateEntityOpt) error

	CreateFeatureTx func(context.Context, *sqlx.Tx, CreateFeatureOpt) (int16, error)
	UpdateFeatureTx func(context.Context, *sqlx.Tx, UpdateFeatureOpt) error

	CreateFeatureGroupTx func(context.Context, *sqlx.Tx, CreateFeatureGroupOpt) (int16, error)
	UpdateFeatureGroupTx func(context.Context, *sqlx.Tx, UpdateFeatureGroupOpt) error

	CreateRevisionTx func(context.Context, *sqlx.Tx, CreateRevisionOpt) (int32, string, error)
	UpdateRevisionTx func(context.Context, *sqlx.Tx, UpdateRevisionOpt) error
}

func (tp *TxProxy) WithTransaction(ctx context.Context, fn func(tx *Tx) error) error {
	tx, err := tp.BeginTxFn(ctx, nil)
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
	return fn(&Tx{tx: tx, tp: tp})
}

func (tp *TxProxy) CreateEntity(ctx context.Context, opt CreateEntityOpt) (id int16, err error) {
	err = tp.WithTransaction(ctx, func(tx *Tx) (err error) {
		id, err = tx.CreateEntity(ctx, opt)
		return err
	})
	return
}

func (tp *TxProxy) UpdateEntity(ctx context.Context, opt UpdateEntityOpt) error {
	return tp.WithTransaction(ctx, func(tx *Tx) error {
		return tx.UpdateEntity(ctx, opt)
	})
}

func (tp *TxProxy) CreateFeature(ctx context.Context, opt CreateFeatureOpt) (id int16, err error) {
	err = tp.WithTransaction(ctx, func(tx *Tx) (err error) {
		id, err = tx.CreateFeature(ctx, opt)
		return err
	})
	return
}

func (tp *TxProxy) UpdateFeature(ctx context.Context, opt UpdateFeatureOpt) error {
	return tp.WithTransaction(ctx, func(tx *Tx) error {
		return tx.UpdateFeature(ctx, opt)
	})
}

func (tp *TxProxy) CreateFeatureGroup(ctx context.Context, opt CreateFeatureGroupOpt) (id int16, err error) {
	err = tp.WithTransaction(ctx, func(tx *Tx) (err error) {
		id, err = tx.CreateFeatureGroup(ctx, opt)
		return err
	})
	return
}

func (tp *TxProxy) UpdateFeatureGroup(ctx context.Context, opt UpdateFeatureGroupOpt) error {
	return tp.WithTransaction(ctx, func(tx *Tx) error {
		return tx.UpdateFeatureGroup(ctx, opt)
	})
}

func (tp *TxProxy) CreateRevision(ctx context.Context, opt CreateRevisionOpt) (id int32, dataTable string, err error) {
	err = tp.WithTransaction(ctx, func(tx *Tx) (err error) {
		id, dataTable, err = tx.CreateRevision(ctx, opt)
		return err
	})
	return id, dataTable, err
}

func (tp *TxProxy) UpdateRevision(ctx context.Context, opt UpdateRevisionOpt) error {
	return tp.WithTransaction(ctx, func(tx *Tx) error {
		return tx.UpdateRevision(ctx, opt)
	})
}
