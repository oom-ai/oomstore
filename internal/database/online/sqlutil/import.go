package sqlutil

import (
	"context"
	"fmt"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
)

func Import(ctx context.Context, db *sqlx.DB, opt online.ImportOpt, backend types.BackendType) error {
	columns := append([]string{opt.Entity.Name}, opt.FeatureList.Names()...)
	err := dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		tableName := OnlineTableName(opt.Revision.ID)
		schema, err := dbutil.BuildCreateSchema(tableName, opt.Entity, opt.FeatureList, backend)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, schema)
		if err != nil {
			return err
		}

		// populate the data table
		records := make([]interface{}, 0, BatchSize)
		for record := range opt.ExportStream {
			if len(record) != len(opt.FeatureList)+1 {
				return fmt.Errorf("field count not matched, expected %d, got %d", len(opt.FeatureList)+1, len(record))
			}
			records = append(records, record)

			if len(records) == BatchSize {
				if err := dbutil.InsertRecordsToTableTx(tx, ctx, tableName, records, columns, backend); err != nil {
					return err
				}
				records = make([]interface{}, 0, BatchSize)
			}
		}

		if err := dbutil.InsertRecordsToTableTx(tx, ctx, tableName, records, columns, backend); err != nil {
			return err
		}
		if opt.ExportError != nil {
			return <-opt.ExportError
		}
		return nil
	})
	return err
}
