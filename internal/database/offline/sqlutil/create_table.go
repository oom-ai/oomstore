package sqlutil

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func CreateTable(ctx context.Context, db *sqlx.DB, opt offline.CreateTableOpt, backend types.BackendType) error {
	if opt.IsCDC {
		// Create index (entity_key, unix_milli) on cdc table
		schema := dbutil.BuildTableSchema(opt.TableName, opt.Entity, true, opt.Features, nil, backend)
		if _, err := db.ExecContext(ctx, schema); err != nil {
			return err
		}
		indexFields := []string{opt.Entity.Name, "unix_milli"}
		indexDDL := dbutil.BuildIndexDDL(opt.TableName, "idx", indexFields, backend)
		if _, err := db.ExecContext(ctx, indexDDL); err != nil {
			return err
		}
	} else {
		// Create primary key (entity_key) on snapshot table
		pkFields := []string{opt.Entity.Name}
		schema := dbutil.BuildTableSchema(opt.TableName, opt.Entity, false, opt.Features, pkFields, backend)
		if _, err := db.ExecContext(ctx, schema); err != nil {
			return err
		}
	}
	return nil
}
