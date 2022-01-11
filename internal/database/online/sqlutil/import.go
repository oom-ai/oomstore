package sqlutil

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const importBatchSize = 10

func Import(ctx context.Context, db *sqlx.DB, opt online.ImportOpt, backend types.BackendType) error {
	columns := append([]string{opt.Entity.Name}, opt.Features.Names()...)
	err := dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		tableName := OnlineBatchTableName(opt.Revision.ID)
		schema := dbutil.BuildTableSchema(tableName, opt.Entity, false, opt.Features, []string{opt.Entity.Name}, backend)
		_, err := tx.ExecContext(ctx, schema)
		if err != nil {
			return errdefs.WithStack(err)
		}

		// populate the data table
		records := make([]interface{}, 0, importBatchSize)
		for record := range opt.ExportStream {
			if len(record) != len(opt.Features)+1 {
				return errdefs.Errorf("field count not matched, expected %d, got %d", len(opt.Features)+1, len(record))
			}
			records = append(records, record)

			if len(records) == importBatchSize {
				if err := dbutil.InsertRecordsToTableTx(tx, ctx, tableName, records, columns, backend); err != nil {
					return err
				}
				records = make([]interface{}, 0, importBatchSize)
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
