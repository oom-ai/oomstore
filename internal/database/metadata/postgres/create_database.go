package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func CreateDatabase(ctx context.Context, opt types.PostgresDbOpt) (err error) {
	db, err := OpenWith(opt.Host, opt.Port, opt.User, opt.Pass, "")
	if err != nil {
		return
	}
	defer db.Close()

	if _, err = db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", opt.Database)); err != nil {
		return
	}

	return createMetaSchemas(ctx, opt)
}

func createMetaSchemas(ctx context.Context, opt types.PostgresDbOpt) (err error) {
	db, err := Open(&opt)
	if err != nil {
		return
	}
	defer db.Close()

	// Use translation to guarantee the following operations be executed
	// on the same connection: http://go-database-sql.org/modifying.html
	return database.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create database functions
		for _, fn := range database.DB_FUNCTIONS {
			if _, err = tx.ExecContext(ctx, fn); err != nil {
				return err
			}
		}

		// create meta tables
		for _, schema := range database.META_TABLE_SCHEMAS {
			if _, err = tx.ExecContext(ctx, schema); err != nil {
				return err
			}
		}

		// create meta views
		for _, schema := range database.META_VIEW_SCHEMAS {
			if _, err = tx.ExecContext(ctx, schema); err != nil {
				return err
			}
		}

		// create triggers
		for table := range database.META_TABLE_SCHEMAS {
			trigger := strings.ReplaceAll(database.TRIGGER_TEMPLATE, `{{TABLE_NAME}}`, table)
			if _, err = tx.ExecContext(ctx, trigger); err != nil {
				return err
			}
		}

		return nil
	})
}
