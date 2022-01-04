package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func CreateDatabase(ctx context.Context, opt *types.PostgresOpt) (err error) {
	defaultDB, err := dbutil.OpenPostgresDB(opt.Host, opt.Port, opt.User, opt.Password, "")
	if err != nil {
		return
	}
	defer defaultDB.Close()

	if _, err = defaultDB.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", opt.Database)); err != nil {
		if e2, ok := err.(*pq.Error); ok && e2.Code != pgerrcode.DuplicateDatabase {
			return err
		}
	}

	db, err := dbutil.OpenPostgresDB(opt.Host, opt.Port, opt.User, opt.Password, opt.Database)
	if err != nil {
		return
	}
	defer db.Close()
	return createMetaSchemas(ctx, db)
}

func createMetaSchemas(ctx context.Context, db *sqlx.DB) (err error) {
	// Use transaction to guarantee the following operations be executed
	// on the same connection: http://go-database-sql.org/modifying.html
	return dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create database functions
		for _, fn := range DB_FUNCTIONS {
			if _, err = tx.ExecContext(ctx, fn); err != nil {
				return err
			}
		}

		// create meta tables
		for _, schema := range META_TABLE_SCHEMAS {
			if _, err = tx.ExecContext(ctx, schema); err != nil {
				return err
			}
		}

		// create foreign keys
		for _, stmt := range META_TABLE_FOREIGN_KEYS {
			if _, err = tx.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}

		// create meta views
		for _, schema := range META_VIEW_SCHEMAS {
			if _, err = tx.ExecContext(ctx, schema); err != nil {
				return err
			}
		}

		// create triggers
		for table := range META_TABLE_SCHEMAS {
			trigger := strings.ReplaceAll(TRIGGER_TEMPLATE, `{{TABLE_NAME}}`, table)
			if _, err = tx.ExecContext(ctx, trigger); err != nil {
				return err
			}
		}

		return nil
	})
}
