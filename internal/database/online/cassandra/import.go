package cassandra

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	entity := opt.Group.Entity
	columns := append([]string{entity.Name}, opt.Features.Names()...)
	tableName := sqlutil.OnlineBatchTableName(opt.Revision.ID)

	table := dbutil.BuildTableSchema(tableName, entity, false, opt.Features, []string{entity.Name}, Backend)

	// create table
	if err := db.Query(table).Exec(); err != nil {
		return errdefs.WithStack(err)
	}

	var (
		insertStmt = buildInsertStatement(tableName, columns)
		batch      = db.NewBatch(gocql.LoggedBatch)
	)
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
