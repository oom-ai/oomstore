package postgres

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

const (
	PostgresBatchSize = 10
)

func (db *DB) SinkFeatureValuesStream(ctx context.Context, stream <-chan *types.RawFeatureValueRecord, features []*types.Feature, revision *types.Revision, entity *types.Entity) error {
	columns := getColumns(*entity, features)
	err := database.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		tmpTableName := revision.GroupName + "_" + strconv.Itoa(rand.Int())
		schema := database.BuildFeatureDataTableSchema(tmpTableName, entity, features)
		_, err := db.ExecContext(ctx, schema)
		if err != nil {
			return err
		}

		// populate the data table
		records := make([]interface{}, 0, PostgresBatchSize)
		for item := range stream {
			if item.Error != nil {
				return item.Error
			}
			record := item.Record
			if len(record) != len(features)+1 {
				return fmt.Errorf("field count not matched, expected %d, got %d", len(features)+1, len(record))
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
		finalTableName := getOnlineBatchTableName(revision)
		rename := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tmpTableName, finalTableName)
		_, err = tx.ExecContext(ctx, rename)
		return err
	})
	return err
}

func getOnlineBatchTableName(revision *types.Revision) string {
	return fmt.Sprintf("batch_%d", revision.ID)
}

func (db *DB) PurgeRevision(ctx context.Context, revision *types.Revision) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, getOnlineBatchTableName(revision))
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

func getColumns(entity types.Entity, features []*types.Feature) []string {
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
