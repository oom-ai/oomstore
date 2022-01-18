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
	err := dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		pkFields := []string{opt.Entity.Name}
		if opt.NoPK {
			pkFields = nil
		}
		schema := dbutil.BuildTableSchema(opt.SnapshotTableName, opt.Entity, false, opt.Features, pkFields, backend)
		_, err := tx.ExecContext(ctx, schema)
		if err != nil {
			return errdefs.WithStack(err)
		}

		// populate the data table
		err = loadData(tx, ctx, dbutil.LoadDataFromSourceOpt{
			Source:    opt.Source,
			Entity:    opt.Entity,
			TableName: opt.SnapshotTableName,
			Header:    opt.Header,
			Features:  opt.Features,
			Backend:   backend,
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
