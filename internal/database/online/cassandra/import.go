package cassandra

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocql/gocql"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	// Step 0: drop existing table for streaming feature
	var tableName string
	if opt.Group.Category == types.CategoryBatch {
		tableName = sqlutil.OnlineBatchTableName(*opt.RevisionID)
	} else {
		tableName = sqlutil.OnlineStreamTableName(opt.Group.ID)
		if err := db.Query(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName)).Exec(); err != nil {
			return errdefs.WithStack(err)
		}
	}

	// Step 1: create online table
	entity := opt.Group.Entity
	columns := append([]string{entity.Name}, opt.Features.Names()...)
	schema := dbutil.BuildTableSchema(tableName, entity.Name, false, opt.Features, []string{entity.Name}, Backend)
	if err := db.Query(schema).Exec(); err != nil {
		return errdefs.WithStack(err)
	}

	// Step 2: insert records to the online table
	insertStmt := buildInsertStatement(tableName, columns)
	batch := db.NewBatch(gocql.LoggedBatch)
	for record := range opt.ExportStream {
		if len(record) != len(opt.Features)+1 {
			return errdefs.Errorf("field count not matched, expected %d, got %d", len(opt.Features)+1, len(record))
		}
		if batch.Size() != BatchSize {
			batch.Query(insertStmt, record...)
		} else {
			if err := db.ExecuteBatch(batch); err != nil {
				return errdefs.WithStack(err)
			}
			batch = db.NewBatch(gocql.LoggedBatch)
		}
	}
	return errdefs.WithStack(db.ExecuteBatch(batch))
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
