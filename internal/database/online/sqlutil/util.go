package sqlutil

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func OnlineBatchTableName(revisionID int) string {
	return fmt.Sprintf("online_batch_%d", revisionID)
}

func OnlineStreamTableName(groupID int) string {
	return fmt.Sprintf("online_stream_%d", groupID)
}

func CreateStreamTableSchema(ctx context.Context, tableName string, entity *types.Entity, backend types.BackendType) (string, error) {
	var (
		entityFormat string
	)
	switch backend {
	case types.BackendMySQL:
		entityFormat = fmt.Sprintf("`%s` VARCHAR(255) PRIMARY KEY", entity.Name)
	case types.BackendPostgres, types.BackendCassandra, types.BackendSQLite:
		entityFormat = fmt.Sprintf(`"%s" TEXT PRIMARY KEY`, entity.Name)
	default:
		return "", errdefs.InvalidAttribute(errors.Errorf("backend %s not support", backend))
	}

	schema := fmt.Sprintf("CREATE TABLE %s ( %s )", tableName, entityFormat)
	return schema, nil
}

func SqlxPrapareStreamTable(ctx context.Context, db *sqlx.DB, opt online.PrepareStreamTableOpt, backend types.BackendType) error {
	tableName := OnlineStreamTableName(opt.GroupID)

	if opt.Feature == nil {
		schema, err := CreateStreamTableSchema(ctx, tableName, opt.Entity, backend)
		if err != nil {
			return err
		}
		_, err = db.ExecContext(ctx, schema)
		return errors.WithStack(err)
	}

	dbValueType, err := dbutil.DBValueType(backend, opt.Feature.ValueType)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, opt.Feature.Name, dbValueType)
	_, err = db.ExecContext(ctx, sql)
	return errors.WithStack(err)
}
