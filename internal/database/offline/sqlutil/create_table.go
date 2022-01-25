package sqlutil

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func CreateTable(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.CreateTableOpt) error {
	if !supportIndex(dbOpt.Backend) {
		return CreateTableNoIndex(ctx, dbOpt, opt)
	}
	switch opt.TableType {
	case types.TableBatchSnapshot:
		// Create primary key (entity_key) on batch snapshot table
		pkFields := []string{opt.EntityName}
		schema := dbutil.BuildTableSchema(opt.TableName, opt.EntityName, false, opt.Features, pkFields, dbOpt.Backend)
		if err := dbOpt.ExecContext(ctx, schema, nil); err != nil {
			return err
		}
	case types.TableStreamSnapshot:
		// Create primary key (entity_key) on stream snapshot table
		pkFields := []string{opt.EntityName}
		schema := dbutil.BuildTableSchema(opt.TableName, opt.EntityName, true, opt.Features, pkFields, dbOpt.Backend)
		if err := dbOpt.ExecContext(ctx, schema, nil); err != nil {
			return err
		}
	case types.TableStreamCdc:
		schema := dbutil.BuildTableSchema(opt.TableName, opt.EntityName, true, opt.Features, nil, dbOpt.Backend)
		if err := dbOpt.ExecContext(ctx, schema, nil); err != nil {
			return err
		}
		// Create index (entity_key, unix_milli) on stream cdc table
		indexFields := []string{opt.EntityName, "unix_milli"}
		indexDDL := dbutil.BuildIndexDDL(opt.TableName, "idx", indexFields, dbOpt.Backend)
		if err := dbOpt.ExecContext(ctx, indexDDL, nil); err != nil {
			return err
		}
	default:
		panic(fmt.Sprintf("unsupported table type %s", opt.TableType))
	}
	return nil
}

func CreateTableNoIndex(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.CreateTableOpt) error {
	var hasUnixMilli bool
	switch opt.TableType {
	case types.TableBatchSnapshot:
		hasUnixMilli = false
	case types.TableStreamSnapshot, types.TableStreamCdc:
		hasUnixMilli = true
	default:
		panic(fmt.Sprintf("unsupported table type %s", opt.TableType))
	}

	tableName := opt.TableName
	if dbOpt.Backend == types.BackendBigQuery {
		tableName = fmt.Sprintf("%s.%s", *dbOpt.DatasetID, tableName)
	}
	schema := dbutil.BuildTableSchema(tableName, opt.EntityName, hasUnixMilli, opt.Features, nil, dbOpt.Backend)
	if err := dbOpt.ExecContext(ctx, schema, nil); err != nil {
		return err
	}
	return nil
}
