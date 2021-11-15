package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	PostgresBatchSize = 10
)

func (db *DB) createTableJoined(ctx context.Context, features types.FeatureList, entity types.Entity, groupName string) (string, error) {
	// create table joined
	tableName := dbutil.TempTable(fmt.Sprintf("joined_%s", groupName))
	schema := `
		CREATE TABLE %s (
			entity_key  VARCHAR(%d) NOT NULL,
			unix_time   BIGINT NOT NULL,
			%s
		);
	`
	index := fmt.Sprintf(`CREATE INDEX ON %s (unix_time, entity_key)`, tableName)

	var columnDefs []string
	for _, f := range features {
		columnDefs = append(columnDefs, fmt.Sprintf(`"%s" %s`, f.Name, f.DBValueType))
	}

	schema = fmt.Sprintf(schema, tableName, entity.Length, strings.Join(columnDefs, ",\n"))
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return "", err
	}
	_, err := db.ExecContext(ctx, index)

	return tableName, err
}

func (db *DB) createAndImportTableEntityRows(ctx context.Context, entity types.Entity, entityRows <-chan types.EntityRow) (string, error) {
	// create table entity_rows
	tableName := dbutil.TempTable("entity_rows")
	schema := fmt.Sprintf(`
		CREATE TABLE %s (
			entity_key  VARCHAR(%d) NOT NULL,
			unix_time   BIGINT NOT NULL
		);
	`, tableName, entity.Length)

	if _, err := db.ExecContext(ctx, schema); err != nil {
		return "", err
	}

	// populate dataset to the table
	if err := db.insertEntityRows(ctx, tableName, entityRows); err != nil {
		return "", err
	}

	// create index
	index := fmt.Sprintf(`CREATE INDEX ON %s (unix_time, entity_key)`, tableName)
	if _, err := db.ExecContext(ctx, index); err != nil {
		return "", err
	}
	return tableName, nil
}

func (db *DB) insertEntityRows(ctx context.Context, tableName string, entityRows <-chan types.EntityRow) error {
	records := make([]interface{}, 0, PostgresBatchSize)
	columns := []string{"entity_key", "unix_time"}
	for entityRow := range entityRows {
		records = append(records, []interface{}{entityRow.EntityKey, entityRow.UnixTime})

		if len(records) == PostgresBatchSize {
			if err := dbutil.InsertRecordsToTable(db.DB, ctx, tableName, records, columns); err != nil {
				return err
			}
			records = make([]interface{}, 0, PostgresBatchSize)
		}
	}
	if err := dbutil.InsertRecordsToTable(db.DB, ctx, tableName, records, columns); err != nil {
		return err
	}
	return nil
}

func (db *DB) dropTable(ctx context.Context, tableName string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName)
	_, err := db.ExecContext(ctx, query)
	return err
}
