package postgres

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/internal/database/dbutil"
	"github.com/onestore-ai/onestore/internal/database/online"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

const (
	PostgresBatchSize = 10
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	columns := getColumns(opt.Entity, opt.Features)
	err := dbutil.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		tmpTableName := opt.Revision.GroupName + "_" + strconv.Itoa(rand.Int())
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
		finalTableName := getOnlineBatchTableName(opt.Revision)
		rename := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tmpTableName, finalTableName)
		_, err = tx.ExecContext(ctx, rename)
		return err
	})
	return err
}

func getOnlineBatchTableName(revision *types.Revision) string {
	return fmt.Sprintf("batch_%d", revision.ID)
}

func (db *DB) Purge(ctx context.Context, revision *types.Revision) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, getOnlineBatchTableName(revision))
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

func getColumns(entity *types.Entity, features []*types.Feature) []string {
	columns := make([]string, 0, len(features)+1)
	columns = append(columns, entity.Name)
	for _, f := range features {
		columns = append(columns, f.Name)
	}
	return columns
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
		fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, strings.Join(columns, ","), strings.Join(valueFlags, ",")),
		records...)
	if err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, db.Rebind(query), args...); err != nil {
		return err
	}
	return nil
}
