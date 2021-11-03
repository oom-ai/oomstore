package postgres

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
)

const (
	PostgresBatchSize = 10
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	columns := append([]string{opt.Entity.Name}, opt.Features.Names()...)
	err := dbutil.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		tmpTableName := fmt.Sprint(opt.Revision.GroupId) + "_" + strconv.Itoa(rand.Int())
		schema := dbutil.BuildFeatureDataTableSchema(tmpTableName, opt.Entity, opt.Features)
		_, err := db.ExecContext(ctx, schema)
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
			if len(record) != len(opt.Features)+1 {
				return fmt.Errorf("field count not matched, expected %d, got %d", len(opt.Features)+1, len(record))
			}
			records = append(records, record)

			if len(records) == PostgresBatchSize {
				if err := db.insertRecordsToTable(ctx, tmpTableName, records, columns); err != nil {
					return err
				}
				records = make([]interface{}, 0, PostgresBatchSize)
			}
		}

		if err := db.insertRecordsToTable(ctx, tmpTableName, records, columns); err != nil {
			return err
		}

		// rename the tmp table to final table
		finalTableName := getOnlineBatchTableName(opt.Revision.ID)
		rename := fmt.Sprintf(`ALTER TABLE "%s" RENAME TO "%s"`, tmpTableName, finalTableName)
		_, err = tx.ExecContext(ctx, rename)
		return err
	})
	return err
}

func getOnlineBatchTableName(revisionId int32) string {
	return fmt.Sprintf("batch_%d", revisionId)
}

func (db *DB) insertRecordsToTable(ctx context.Context, tableName string, records []interface{}, columns []string) error {
	if len(records) == 0 {
		return nil
	}
	valueFlags := make([]string, 0, PostgresBatchSize)
	for i := 0; i < len(records); i++ {
		valueFlags = append(valueFlags, "(?)")
	}

	query, args, err := sqlx.In(
		fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES %s`, tableName, dbutil.Quote(`"`, columns...), strings.Join(valueFlags, ",")),
		records...)
	if err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, db.Rebind(query), args...); err != nil {
		return err
	}
	return nil
}
