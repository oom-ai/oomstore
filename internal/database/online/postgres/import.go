package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
)

const (
	PostgresBatchSize = 10
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	columns := append([]string{opt.Entity.Name}, opt.FeatureList.Names()...)
	err := dbutil.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		tmpTableName := dbutil.TempTable(fmt.Sprintf("import_revision_%d", opt.Revision.ID))
		schema := dbutil.BuildFeatureDataTableSchema(tmpTableName, opt.Entity, opt.FeatureList)
		_, err := tx.ExecContext(ctx, schema)
		if err != nil {
			return err
		}

		// populate the data table
		records := make([]interface{}, 0, PostgresBatchSize)
		for item := range opt.Stream {
			if item.Error != nil {
				return item.Error
			}
			record := item.Record
			if len(record) != len(opt.FeatureList)+1 {
				return fmt.Errorf("field count not matched, expected %d, got %d", len(opt.FeatureList)+1, len(record))
			}
			records = append(records, record)

			if len(records) == PostgresBatchSize {
				if err := dbutil.InsertRecordsToTableTx(tx, ctx, tmpTableName, records, columns); err != nil {
					return err
				}
				records = make([]interface{}, 0, PostgresBatchSize)
			}
		}

		if err := dbutil.InsertRecordsToTableTx(tx, ctx, tmpTableName, records, columns); err != nil {
			return err
		}

		// rename the tmp table to final table
		finalTableName := getOnlineTableName(opt.Revision.ID)
		rename := fmt.Sprintf(`ALTER TABLE "%s" RENAME TO "%s"`, tmpTableName, finalTableName)
		_, err = tx.ExecContext(ctx, rename)
		return err
	})
	return err
}

func getOnlineTableName(revisionID int) string {
	return fmt.Sprintf("online_%d", revisionID)
}
