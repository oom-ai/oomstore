package database

import (
	"context"
)

func CreateDatabase(ctx context.Context, opt Option) (err error) {
	db, err := OpenWith(opt.Host, opt.Port, opt.User, opt.Pass, "")
	if err != nil {
		return
	}
	defer db.Close()
	if _, err = db.ExecContext(ctx, "CREATE DATABASE ?", opt.DbName); err != nil {
		return
	}

	// Use translation to guarantee the following operations be executed
	// on the same connection: http://go-database-sql.org/modifying.html
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	if _, err = tx.ExecContext(ctx, "USE ?", opt.DbName); err != nil {
		return
	}

	for _, schema := range META_SCHEMAS {
		if _, err = tx.ExecContext(ctx, schema); err != nil {
			return
		}
	}

	return
}
