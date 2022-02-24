package sqlutil

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const importBatchSize = 10

func Import(ctx context.Context, db *sqlx.DB, opt online.ImportOpt, backend types.BackendType) error {
	var tableName string
	if opt.Group.Category == types.CategoryStream {
		tableName = dbutil.TempTable("online_stream")
	} else {
		tableName = dbutil.OnlineBatchTableName(*opt.RevisionID)
	}
	entity := opt.Group.Entity
	columns := append([]string{entity.Name}, opt.Features.Names()...)

	return dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		schema := dbutil.BuildTableSchema(dbutil.BuildTableSchemaParams{
			TableName:    tableName,
			EntityName:   entity.Name,
			HasUnixMilli: false,
			Features:     opt.Features,
			PrimaryKeys:  []string{entity.Name},
			Backend:      backend,
		})
		_, err := tx.ExecContext(ctx, schema)
		if err != nil {
			return errdefs.WithStack(err)
		}

		// populate the data table
		records := make([]interface{}, 0, importBatchSize)
		for record := range opt.ExportStream {
			if record.Error != nil {
				return record.Error
			}

			if len(record.Record) != len(opt.Features)+1 {
				return errdefs.Errorf("field count not matched, expected %d, got %d", len(opt.Features)+1, len(record.Record))
			}
			records = append(records, record.Record)

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

		if opt.Group.Category == types.CategoryStream {
			streamTableName := dbutil.OnlineStreamTableName(opt.Group.ID)
			if err := PurgeTx(ctx, tx, streamTableName, backend); err != nil {
				return err
			}
			query := fmt.Sprintf(`ALTER TABLE %s RENAME TO %s;`, tableName, streamTableName)
			if _, err = tx.ExecContext(ctx, query); err != nil {
				return errdefs.WithStack(err)
			}
		}
		return nil
	})
}
