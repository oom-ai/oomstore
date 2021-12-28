package sqlutil

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type LoadData func(tx *sqlx.Tx, ctx context.Context, source *offline.CSVSource, tableName string, header []string) error

func Import(ctx context.Context, db *sqlx.DB, opt offline.ImportOpt, loadData LoadData, backendType types.BackendType) (int64, error) {
	var revision int64
	err := dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		schema, err := dbutil.BuildCreateSchema(opt.SnapshotTableName, opt.Entity, opt.Features, backendType)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, schema)
		if err != nil {
			return err
		}

		// populate the data table
		err = loadData(tx, ctx, opt.Source, opt.SnapshotTableName, opt.Header)
		if err != nil {
			return err
		}

		if opt.Revision != nil {
			// use user-defined revision
			revision = *opt.Revision
		} else {
			// generate revision using current timestamp
			revision = time.Now().UnixMilli()
		}
		return nil
	})
	return revision, err
}
