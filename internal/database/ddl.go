package database

import (
	"context"
)

func CreateDatabase(ctx context.Context, opt Option) error {
	db, err := OpenWith(opt.Host, opt.Port, opt.User, opt.Pass, "")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.ExecContext(ctx, "CREATE DATABASE ?", opt.DbName)
	if err != nil {
		return err
	}

	// Use translation to guarantee the following operations be executed
	// on the same connection: http://go-database-sql.org/modifying.html
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "USE ?", opt.DbName); err != nil {
		return err
	}

	for _, schema := range META_SCHEMAS {
		if _, err := tx.ExecContext(ctx, schema); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
