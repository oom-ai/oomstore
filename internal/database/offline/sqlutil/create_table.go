package sqlutil

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func CreateTable(ctx context.Context, db *sqlx.DB, opt offline.CreateTableOpt, backend types.BackendType) error {
	switch opt.TableType {
	case types.TableBatchSnapshot:
		// Create primary key (entity_key) on batch snapshot table
		pkFields := []string{opt.Entity.Name}
		schema := dbutil.BuildTableSchema(opt.TableName, opt.Entity, false, opt.Features, pkFields, backend)
		if _, err := db.ExecContext(ctx, schema); err != nil {
			return errdefs.WithStack(err)
		}
	case types.TableStreamSnapshot:
		// Create primary key (entity_key) on stream snapshot table
		pkFields := []string{opt.Entity.Name}
		schema := dbutil.BuildTableSchema(opt.TableName, opt.Entity, true, opt.Features, pkFields, backend)
		if _, err := db.ExecContext(ctx, schema); err != nil {
			return errdefs.WithStack(err)
		}
	case types.TableStreamCdc:
		schema := dbutil.BuildTableSchema(opt.TableName, opt.Entity, true, opt.Features, nil, backend)
		if _, err := db.ExecContext(ctx, schema); err != nil {
			return errdefs.WithStack(err)
		}
		// Create index (entity_key, unix_milli) on stream cdc table
		indexFields := []string{opt.Entity.Name, "unix_milli"}
		indexDDL := dbutil.BuildIndexDDL(opt.TableName, "idx", indexFields, backend)
		if _, err := db.ExecContext(ctx, indexDDL); err != nil {
			return errdefs.WithStack(err)
		}
	default:
		panic(fmt.Sprintf("unsupported table type %s", opt.TableType))
	}
	return nil
}
