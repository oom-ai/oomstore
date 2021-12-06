package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

const (
	PostgresBatchSize = 10
)

func (db *DB) createTableJoined(ctx context.Context, features types.FeatureList, entity types.Entity, groupName string, valueNames []string) (string, error) {
	// create table joined
	tableName := dbutil.TempTable(fmt.Sprintf("joined_%s", groupName))
	columnDefs := []string{fmt.Sprintf("entity_key  VARCHAR(%d) NOT NULL", entity.Length), "unix_milli  BIGINT NOT NULL"}
	for _, name := range valueNames {
		columnDefs = append(columnDefs, fmt.Sprintf(`"%s" TEXT`, name))
	}
	for _, f := range features {
		columnDefs = append(columnDefs, fmt.Sprintf(`"%s" %s`, f.Name, f.DBValueType))
	}
	schema := `
		CREATE TABLE %s (
			%s
		);
	`
	index := fmt.Sprintf(`CREATE INDEX ON %s (unix_milli, entity_key)`, tableName)

	schema = fmt.Sprintf(schema, tableName, strings.Join(columnDefs, ",\n"))
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return "", err
	}
	_, err := db.ExecContext(ctx, index)

	return tableName, err
}

func (db *DB) createAndImportTableEntityRows(ctx context.Context, entity types.Entity, entityRows <-chan types.EntityRow, valueNames []string) (string, error) {
	// create table entity_rows
	tableName := dbutil.TempTable("entity_rows")
	columnDefs := []string{fmt.Sprintf("entity_key  VARCHAR(%d) NOT NULL", entity.Length), "unix_milli  BIGINT NOT NULL"}
	for _, name := range valueNames {
		columnDefs = append(columnDefs, fmt.Sprintf(`"%s" %s`, name, "TEXT"))
	}
	schema := fmt.Sprintf(`
		CREATE TABLE %s (
			%s
		);
	`, tableName, strings.Join(columnDefs, ",\n"))

	if _, err := db.ExecContext(ctx, schema); err != nil {
		return "", err
	}

	// populate dataset to the table
	if err := db.insertEntityRows(ctx, tableName, entityRows, valueNames); err != nil {
		return "", err
	}

	// create index
	index := fmt.Sprintf(`CREATE INDEX ON %s (unix_milli, entity_key)`, tableName)
	if _, err := db.ExecContext(ctx, index); err != nil {
		return "", err
	}
	return tableName, nil
}

func (db *DB) insertEntityRows(ctx context.Context, tableName string, entityRows <-chan types.EntityRow, valueNames []string) error {
	records := make([]interface{}, 0, PostgresBatchSize)
	columns := []string{"entity_key", "unix_milli"}
	columns = append(columns, valueNames...)
	for entityRow := range entityRows {
		record := []interface{}{entityRow.EntityKey, entityRow.UnixMilli}
		for _, v := range entityRow.Values {
			record = append(record, v)
		}
		records = append(records, record)
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
