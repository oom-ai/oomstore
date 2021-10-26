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
	valueFlags := make([]string, 0, PostgresBatchSize)
	for i := 0; i < PostgresBatchSize; i++ {
		valueFlags = append(valueFlags, "(?)")
	}
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
				query, args, err := sqlx.In(
					fmt.Sprintf("INSERT INTO %s (%s) VALUES IN %s", tmpTableName, strings.Join(columns, ","), strings.Join(valueFlags, ",")),
					records...)
				if err != nil {
					return err
				}
				if _, err := db.Exec(query, args...); err != nil {
					return err
				}
				records = make([]interface{}, 0, PostgresBatchSize)
			}
		}

		// rename the tmp table to final table
		finalTableName := revision.GroupName
		query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, finalTableName)
		if _, err := db.ExecContext(ctx, query); err != nil {
			return err
		}
		rename := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tmpTableName, finalTableName)
		_, err = tx.ExecContext(ctx, rename)
		return err
	})
	return err
}

func getColumns(entity types.Entity, features []*types.Feature) []string {
	columns := make([]string, 0, len(features)+1)
	columns = append(columns, entity.Name)
	for _, f := range features {
		columns = append(columns, f.Name)
	}
	return columns
}
