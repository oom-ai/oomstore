package sqlutil

import (
	"context"
	"time"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type LoadData func(tx *sqlx.Tx, ctx context.Context, opt dbutil.LoadDataFromSourceOpt) error

func Import(ctx context.Context, db *sqlx.DB, opt offline.ImportOpt, loadData LoadData, backend types.BackendType) (int64, error) {
	var revision int64
	params := dbutil.BuildTableSchemaParams{
		TableName:    opt.SnapshotTableName,
		EntityName:   opt.EntityName,
		HasUnixMilli: false,
		Features:     opt.Features,
		Backend:      backend,
	}
	err := dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table: streaming snapshot table doesn't have primary key
		if opt.Category == types.CategoryBatch {
			params.PrimaryKeys = []string{opt.EntityName}
		}
		schema := dbutil.BuildTableSchema(params)
		_, err := tx.ExecContext(ctx, schema)
		if err != nil {
			return errdefs.WithStack(err)
		}

		// populate the data table
		err = loadData(tx, ctx, dbutil.LoadDataFromSourceOpt{
			Source:     opt.Source,
			EntityName: opt.EntityName,
			TableName:  opt.SnapshotTableName,
			Header:     opt.Header,
			Features:   opt.Features,
			Backend:    backend,
		})
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
