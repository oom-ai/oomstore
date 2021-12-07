package sqlutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
)

const (
	InsertBatchSize = 20
)

func createTableJoined(ctx context.Context, db *sqlx.DB, features types.FeatureList, entity types.Entity, groupName string, valueNames []string, backendType types.BackendType) (string, error) {
	columnFormat, err := dbutil.GetColumnFormat(backendType)
	if err != nil {
		return "", err
	}
	// create table joined
	tableName := dbutil.TempTable(fmt.Sprintf("joined_%s", groupName))
	columnDefs := []string{fmt.Sprintf("entity_key  VARCHAR(%d) NOT NULL", entity.Length), "unix_milli  BIGINT NOT NULL"}
	for _, name := range valueNames {
		columnDefs = append(columnDefs, fmt.Sprintf(columnFormat, name, "TEXT"))
	}
	for _, f := range features {
		columnDefs = append(columnDefs, fmt.Sprintf(columnFormat, f.Name, f.DBValueType))
	}
	schema := `
		CREATE TABLE %s (
			%s
		);
	`
	index := fmt.Sprintf(`CREATE INDEX idx_%s ON %s (unix_milli, entity_key)`, tableName, tableName)

	schema = fmt.Sprintf(schema, tableName, strings.Join(columnDefs, ",\n"))
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return "", err
	}
	_, err = db.ExecContext(ctx, index)

	return tableName, err
}

func createAndImportTableEntityRows(ctx context.Context, db *sqlx.DB, entity types.Entity, entityRows <-chan types.EntityRow, valueNames []string, backendType types.BackendType) (string, error) {
	columnFormat, err := dbutil.GetColumnFormat(backendType)
	if err != nil {
		return "", err
	}

	// create table entity_rows
	tableName := dbutil.TempTable("entity_rows")
	columnDefs := []string{fmt.Sprintf("entity_key  VARCHAR(%d) NOT NULL", entity.Length), "unix_milli  BIGINT NOT NULL"}
	for _, name := range valueNames {
		columnDefs = append(columnDefs, fmt.Sprintf(columnFormat, name, "TEXT"))
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
	if err := insertEntityRows(ctx, db, tableName, entityRows, valueNames, backendType); err != nil {
		return "", err
	}

	// create index
	index := fmt.Sprintf(`CREATE INDEX idx_%s ON %s (unix_milli, entity_key)`, tableName, tableName)
	if _, err := db.ExecContext(ctx, index); err != nil {
		return "", err
	}
	return tableName, nil
}

func insertEntityRows(ctx context.Context, db *sqlx.DB, tableName string, entityRows <-chan types.EntityRow, valueNames []string, backendType types.BackendType) error {
	records := make([]interface{}, 0, InsertBatchSize)
	columns := []string{"entity_key", "unix_milli"}
	columns = append(columns, valueNames...)
	for entityRow := range entityRows {
		record := []interface{}{entityRow.EntityKey, entityRow.UnixMilli}
		for _, v := range entityRow.Values {
			record = append(record, v)
		}
		records = append(records, record)
		if len(records) == InsertBatchSize {
			if err := dbutil.InsertRecordsToTable(db, ctx, tableName, records, columns, backendType); err != nil {
				return err
			}
			records = make([]interface{}, 0, InsertBatchSize)
		}
	}
	if err := dbutil.InsertRecordsToTable(db, ctx, tableName, records, columns, backendType); err != nil {
		return err
	}
	return nil
}

func dropTable(ctx context.Context, db *sqlx.DB, tableName string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName)
	_, err := db.ExecContext(ctx, query)
	return err
}
