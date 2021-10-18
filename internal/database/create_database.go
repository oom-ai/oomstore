package database

import (
	"context"
	"database/sql"
	"fmt"
)

func CreateDatabase(ctx context.Context, opt Option) (err error) {
	db, err := OpenWith(opt.Host, opt.Port, opt.User, opt.Pass, "")
	if err != nil {
		return
	}
	defer db.Close()
	if _, err = db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE `%s`", opt.DbName)); err != nil {
		return
	}

	// Use translation to guarantee the following operations be executed
	// on the same connection: http://go-database-sql.org/modifying.html
	return db.WithTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if _, err = tx.ExecContext(ctx, fmt.Sprintf("USE `%s`", opt.DbName)); err != nil {
			return err
		}
		// create meta tables
		for _, schema := range META_TABLE_SCHEMAS {
			if _, err = tx.ExecContext(ctx, schema); err != nil {
				return err
			}
		}

		// create meta views
		for _, schema := range META_VIEW_SCHEMAS {
			if _, err = tx.ExecContext(ctx, schema); err != nil {
				return err
			}
		}

		return nil
	})
}
