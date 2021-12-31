package cassandra

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocql/gocql"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	columns := append([]string{opt.Entity.Name}, opt.FeatureList.Names()...)
	tableName := sqlutil.OnlineBatchTableName(opt.Revision.ID)

	table := dbutil.BuildTableSchema(tableName, opt.Entity, false, opt.FeatureList, Backend)

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
			if err := db.ExecuteBatch(batch); err != nil {
				return err
			}
			batch = db.NewBatch(gocql.LoggedBatch)
		}
	}
	return db.ExecuteBatch(batch)
}

func buildInsertStatement(tableName string, columns []string) string {
	valueFlags := make([]string, 0, len(columns))
	for i := 0; i < len(columns); i++ {
		valueFlags = append(valueFlags, "?")
	}

	return fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`,
		tableName,
		dbutil.QuoteFn(Backend)(columns...),
		strings.Join(valueFlags, ","))
}
