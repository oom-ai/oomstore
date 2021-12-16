package cassandra

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocql/gocql"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/internal/database/online/sqlutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
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

func getDbTypeFrom(valueType string) (string, error) {
	if t, ok := typeMap[valueType]; !ok {
		return "", fmt.Errorf("unsupported value type: %s", valueType)
	} else {
		return t, nil
	}
}

func buildDataTableSchema(tableName string, entity *types.Entity, features types.FeatureList) (string, error) {
	columns := make([]dbutil.Column, 0, len(features))
	for _, feature := range features {
		dbType, err := getDbTypeFrom(feature.ValueType)
		if err != nil {
			return "", err
		}

		columns = append(columns, dbutil.Column{
			Name:   feature.Name,
			DbType: dbType,
		})
	}

	return dbutil.BuildSchema(dbutil.Schema{
		TableName:  tableName,
		EntityName: entity.Name,
		Columns:    columns,
	}, dbutil.Cassandra)
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
