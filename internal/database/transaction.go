package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Reference: https://pseudomuto.com/2018/01/clean-sql-transactions-in-golang/

// A Txfn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(ctx context.Context, tx *sqlx.Tx) error

func WithTransaction(db *sqlx.DB, ctx context.Context, fn TxFn) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
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

	return fn(ctx, tx)
}
