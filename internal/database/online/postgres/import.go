package postgres

import (
	"context"
	"fmt"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/jmoiron/sqlx"
)

const (
	PostgresBatchSize = 10
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	columns := append([]string{opt.Entity.Name}, opt.FeatureList.Names()...)
	err := dbutil.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		tableName := getOnlineTableName(opt.Revision.ID)
		schema := dbutil.BuildFeatureDataTableSchema(tableName, opt.Entity, opt.FeatureList)
		_, err := tx.ExecContext(ctx, schema)
		if err != nil {
			return err
		}

		// populate the data table
		records := make([]interface{}, 0, PostgresBatchSize)
		for record := range opt.ExportStream {
			if len(record) != len(opt.FeatureList)+1 {
				return fmt.Errorf("field count not matched, expected %d, got %d", len(opt.FeatureList)+1, len(record))
			}
			records = append(records, record)

			if len(records) == PostgresBatchSize {
				if err := dbutil.InsertRecordsToTableTx(tx, ctx, tableName, records, columns); err != nil {
					return err
				}
				records = make([]interface{}, 0, PostgresBatchSize)
			}
		}

		if err := dbutil.InsertRecordsToTableTx(tx, ctx, tableName, records, columns); err != nil {
			return err
		}
		if opt.ExportError != nil {
			return <-opt.ExportError
		}
		return nil
	})
	return err
}

func getOnlineTableName(revisionID int) string {
	return fmt.Sprintf("online_%d", revisionID)
}
