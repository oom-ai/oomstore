package txproxy

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/metadata"
)

type Tx struct {
	*sqlx.Tx
	*TxProxy
}

func (tx *Tx) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int16, error) {
	return tx.CreateEntityTx(ctx, tx.Tx, opt)
}
func (tx *Tx) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	return tx.UpdateEntityTx(ctx, tx.Tx, opt)
}
func (tx *Tx) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int16, error) {
	return tx.CreateFeatureTx(ctx, tx.Tx, opt)
}
func (tx *Tx) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	return tx.UpdateFeatureTx(ctx, tx.Tx, opt)
}
func (tx *Tx) CreateFeatureGroup(ctx context.Context, opt metadata.CreateFeatureGroupOpt) (int16, error) {
	return tx.CreateFeatureGroupTx(ctx, tx.Tx, opt)
}
func (tx *Tx) UpdateFeatureGroup(ctx context.Context, opt metadata.UpdateFeatureGroupOpt) error {
	return tx.UpdateFeatureGroupTx(ctx, tx.Tx, opt)
}
func (tx *Tx) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int32, string, error) {
	return tx.CreateRevisionTx(ctx, tx.Tx, opt)
}
func (tx *Tx) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	return tx.UpdateRevisionTx(ctx, tx.Tx, opt)
}

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

func (tp *TxProxy) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int16, error) {
	var id int16
	err := tp.WithTransaction(ctx, func(tx *Tx) (err error) {
		id, err = tx.CreateEntity(ctx, opt)
		return err
	})
	return id, err
}

func (tp *TxProxy) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	return tp.WithTransaction(ctx, func(tx *Tx) error {
		return tx.UpdateEntity(ctx, opt)
	})
}

func (tp *TxProxy) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int16, error) {
	var id int16
	err := tp.WithTransaction(ctx, func(tx *Tx) (err error) {
		id, err = tx.CreateFeature(ctx, opt)
		return err
	})
	return id, err
}

func (tp *TxProxy) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	return tp.WithTransaction(ctx, func(tx *Tx) error {
		return tx.UpdateFeature(ctx, opt)
	})
}

func (tp *TxProxy) CreateFeatureGroup(ctx context.Context, opt metadata.CreateFeatureGroupOpt) (int16, error) {
	var id int16
	err := tp.WithTransaction(ctx, func(tx *Tx) (err error) {
		id, err = tx.CreateFeatureGroup(ctx, opt)
		return err
	})
	return id, err
}

func (tp *TxProxy) UpdateFeatureGroup(ctx context.Context, opt metadata.UpdateFeatureGroupOpt) error {
	return tp.WithTransaction(ctx, func(tx *Tx) error {
		return tx.UpdateFeatureGroup(ctx, opt)
	})
}

func (tp *TxProxy) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int32, string, error) {
	var id int32
	var dataTable string
	err := tp.WithTransaction(ctx, func(tx *Tx) (err error) {
		id, dataTable, err = tx.CreateRevision(ctx, opt)
		return err
	})
	return id, dataTable, err
}

func (tp *TxProxy) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	return tp.WithTransaction(ctx, func(tx *Tx) error {
		return tx.UpdateRevision(ctx, opt)
	})
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
	return fn(&Tx{Tx: tx, TxProxy: tp})
}
