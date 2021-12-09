package cassandra

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocql/gocql"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/internal/database/online/sqlutil"
)

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	panic("implement me!")
}

func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	panic("implement me!")
}

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	columns := append([]string{opt.Entity.Name}, opt.FeatureList.Names()...)
	tableName := sqlutil.OnlineTableName(opt.Revision.ID)

	table, err := dbutil.BuildFeatureDataTableSchema(tableName, opt.Entity, opt.FeatureList, BackendType)
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

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	panic("implement me!")
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
