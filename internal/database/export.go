package database

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) createTableEntityDfWithFeatures(ctx context.Context, features []*types.RichFeature, entityName string) (string, error) {
	entity, err := db.GetEntity(ctx, entityName)
	if err != nil {
		return "", err
	}
	tableName := fmt.Sprintf("entity_df_with_features_%d", rand.Int())
	schema := `
		CREATE TABLE %s (
			unique_key  VARCHAR(%d) NOT NULL,
			entity_key  VARCHAR(%d) NOT NULL,
			unix_time   BIGINT NOT NULL,
			%s,
			PRIMARY KEY (entity_key, unix_time)
		);
	`

	var columnDefs []string
	for _, f := range features {
		columnDefs = append(columnDefs, fmt.Sprintf(`"%s" %s`, f.Name, f.ValueType))
	}
	// unique_key = entity_key,unix_time; length of unique_key = entity.Length + 9
	schema = fmt.Sprintf(schema, tableName, entity.Length+10, entity.Length, strings.Join(columnDefs, ",\n"))
	_, err = db.ExecContext(ctx, schema)
	return tableName, err
}

func (db *DB) createAndImportTableEntityDf(ctx context.Context, entityRows []types.EntityRow, entityName string) (string, error) {
	entity, err := db.GetEntity(ctx, entityName)
	if err != nil {
		return "", err
	}
	tableName := fmt.Sprintf("entity_df_%d", rand.Int())
	schema := fmt.Sprintf(`
		CREATE TABLE %s (
			entity_key  VARCHAR(%d) NOT NULL,
			unix_time   BIGINT NOT NULL,
			PRIMARY KEY (entity_key, unix_time)
		);
	`, tableName, entity.Length)
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return tableName, err
	}

	insertQuery := fmt.Sprintf(`INSERT INTO %s(entity_key, unix_time) VALUES (:entity_key, :unix_time)`, tableName)
	_, err = db.NamedExec(insertQuery, entityRows)
	return tableName, err
}

func (db *DB) dropTable(ctx context.Context, tableName string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName)
	_, err := db.ExecContext(ctx, query)
	return err
}
