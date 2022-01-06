package sqlutil

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type LoadData func(tx *sqlx.Tx, ctx context.Context, source *offline.CSVSource, tableName string, header []string, features types.FeatureList) error

func Import(ctx context.Context, db *sqlx.DB, opt offline.ImportOpt, loadData LoadData, backendType types.BackendType) (int64, error) {
	var revision int64
	err := dbutil.WithTransaction(db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		pkFields := []string{opt.Entity.Name}
		if opt.NoPK {
			pkFields = nil
		}
		schema := dbutil.BuildTableSchema(opt.SnapshotTableName, opt.Entity, false, opt.Features, pkFields, backendType)
		_, err := tx.ExecContext(ctx, schema)
		if err != nil {
			return errors.WithStack(err)
		}

		// populate the data table
		err = loadData(tx, ctx, opt.Source, opt.SnapshotTableName, opt.Header, opt.Features)
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
