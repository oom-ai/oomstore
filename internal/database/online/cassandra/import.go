package cassandra

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gocql/gocql"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/internal/database/online/sqlutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

const (
	CASSANDRA_SCHEMA = `CREATE TABLE {{TABLE_NAME}} (
	{{ENTITY_NAME}} TEXT PRIMARY KEY,
	{{COLUMN_DEFS}});
`
	columnFormat = "%s %s"
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	columns := append([]string{opt.Entity.Name}, opt.FeatureList.Names()...)
	tableName := sqlutil.OnlineTableName(opt.Revision.ID)

	table, err := buildDataTableSchema(tableName, opt.Entity, opt.FeatureList)
	if err != nil {
		return err
	}

	// create table
	if err := db.Query(table).Exec(); err != nil {
		return err
	}

	var (
		insertStmt = buildInsertStatement(tableName, columns)
		batch      = db.NewBatch(gocql.LoggedBatch)
	)
	for record := range opt.ExportStream {
		if len(record) != len(opt.FeatureList)+1 {
			return fmt.Errorf("field count not matched, expected %d, got %d", len(opt.FeatureList)+1, len(record))
		}

		if batch.Size() != BatchSize {
			batch.Query(insertStmt, record...)
		} else {
			if err = db.ExecuteBatch(batch); err != nil {
				return err
			}
			batch = db.NewBatch(gocql.LoggedBatch)
		}
	}
	return db.ExecuteBatch(batch)
}

func buildDataTableSchema(tableName string, entity *types.Entity, features types.FeatureList) (string, error) {
	var columnDefs []string
	for _, column := range features {
		dbValueType, err := getDbTypeFrom(column.ValueType)
		if err != nil {
			return "", err
		}
		def := fmt.Sprintf(columnFormat, column.Name, dbValueType)
		columnDefs = append(columnDefs, def)
	}

	// fill schema template
	schema := strings.ReplaceAll(CASSANDRA_SCHEMA, "{{TABLE_NAME}}", tableName)
	schema = strings.ReplaceAll(schema, "{{ENTITY_NAME}}", entity.Name)
	schema = strings.ReplaceAll(schema, "{{ENTITY_LENGTH}}", strconv.Itoa(entity.Length))
	schema = strings.ReplaceAll(schema, "{{COLUMN_DEFS}}", strings.Join(columnDefs, ",\n"))

	return schema, nil
}

var (
	typeMap = map[string]string{
		types.STRING:  "text",
		types.INT64:   "bigint",
		types.FLOAT64: "double",
		types.BOOL:    "boolean",
		types.TIME:    "timestamp",
		types.BYTES:   "text",
	}
)

func getDbTypeFrom(valueType string) (string, error) {
	if t, ok := typeMap[valueType]; !ok {
		return "", fmt.Errorf("unsupported value type: %s", valueType)
	} else {
		return t, nil
	}
}

func buildInsertStatement(tableName string, columns []string) string {
	valueFlags := make([]string, 0, len(columns))
	for i := 0; i < len(columns); i++ {
		valueFlags = append(valueFlags, "?")
	}

	return fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`,
		tableName,
		dbutil.Quote(`"`, columns...),
		strings.Join(valueFlags, ","))
}
