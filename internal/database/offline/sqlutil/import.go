package sqlutil

import (
	"context"
	"encoding/csv"
	"time"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
)

type LoadData func(tx *sqlx.Tx, ctx context.Context, csvReader *csv.Reader, tableName string, header []string) error

func Import(ctx context.Context, db *sqlx.DB, opt offline.ImportOpt, loadData LoadData, backendType types.BackendType) (int64, error) {
	var revision int64
	err := dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		schema, err := dbutil.BuildFeatureDataTableSchema(opt.DataTableName, opt.Entity, opt.Features, backendType)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, schema)
		if err != nil {
			return err
		}

		// populate the data table
		err = loadData(tx, ctx, opt.CsvReader, opt.DataTableName, opt.Header)
		if err != nil {
			return err
		}

		if opt.Revision != nil {
			// use user-defined revision
			revision = *opt.Revision
		} else {
			// generate revision using current timestamp
			revision = time.Now().UnixMilli()
		}
		return nil
	})
	return revision, err
}
