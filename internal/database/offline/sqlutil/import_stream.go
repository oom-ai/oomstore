package sqlutil

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/offline"
)

const (
	SQLBatchSize = 20
)

func ImportStream(ctx context.Context, db *sqlx.DB, opt offline.ImportStreamOpt, backend types.BackendType) (*offline.TimeRange, error) {
	dbOpt := dbutil.DBOpt{
		SqlxDB:  db,
		Backend: backend,
	}
	doImportStreamOpt := DoImportStreamOpt{
		ImportStreamOpt: opt,
		LoadData:        dbutil.LoadDataFromSource(backend, SQLBatchSize),
	}
	return DoImportStream(ctx, dbOpt, doImportStreamOpt)
}

type DoImportStreamOpt struct {
	offline.ImportStreamOpt
	LoadData LoadData
}

func DoImportStream(ctx context.Context, dbOpt dbutil.DBOpt, opt DoImportStreamOpt) (*offline.TimeRange, error) {
	var timeRange offline.TimeRange
	err := dbutil.WithTransaction(dbOpt.SqlxDB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// Step 1: create cdc table
		tempTableName := dbutil.TempTable("offline_stream_cdc")
		schema := dbutil.BuildTableSchema(tempTableName, opt.Entity, true, opt.Features, nil, dbOpt.Backend)
		_, err := tx.ExecContext(ctx, schema)
		if err != nil {
			return errdefs.WithStack(err)
		}
		// Step 2: load data to the cdc table
		err = opt.LoadData(tx, ctx, dbutil.LoadDataFromSourceOpt{
			Source:    opt.Source,
			Entity:    opt.Entity,
			TableName: tempTableName,
			Header:    opt.Header,
			Features:  opt.Features,
			Backend:   dbOpt.Backend,
		})
		if err != nil {
			return err
		}

		// Step 3: find time range of the cdc table, generate cdc and snapshot tables
		res, err := getCdcTimeRange(ctx, tx, tempTableName, dbOpt.Backend)
		if err != nil {
			return err
		}
		timeRange = *res
		group := opt.Features[0].Group
		cdcTableName := dbutil.OfflineStreamCdcTableName(group.ID, timeRange.MinUnixMilli)
		qt := dbutil.QuoteFn(dbOpt.Backend)
		if _, err = tx.ExecContext(ctx, fmt.Sprintf("ALTER TABLE %s RENAME TO %s", qt(tempTableName), qt(cdcTableName))); err != nil {
			return errdefs.WithStack(err)
		}
		if err = CreateTable(ctx, dbOpt.SqlxDB, offline.CreateTableOpt{
			TableName: dbutil.OfflineStreamSnapshotTableName(group.ID, timeRange.MinUnixMilli),
			Entity:    opt.Entity,
			Features:  opt.Features,
			TableType: types.TableStreamSnapshot,
		}, dbOpt.Backend); err != nil {
			return err
		}
		if err = CreateTable(ctx, dbOpt.SqlxDB, offline.CreateTableOpt{
			TableName: dbutil.OfflineStreamCdcTableName(group.ID, timeRange.MaxUnixMilli),
			Entity:    opt.Entity,
			Features:  opt.Features,
			TableType: types.TableStreamCdc,
		}, dbOpt.Backend); err != nil {
			return err
		}

		// Step 4: generate second revision's snapshot table
		if err = Snapshot(ctx, dbOpt, offline.SnapshotOpt{
			Group:        group,
			Features:     opt.Features,
			Revision:     timeRange.MaxUnixMilli,
			PrevRevision: timeRange.MinUnixMilli,
		}); err != nil {
			return err
		}

		return nil
	})
	return &timeRange, err
}

func getCdcTimeRange(ctx context.Context, tx *sqlx.Tx, tableName string, backend types.BackendType) (*offline.TimeRange, error) {
	qt := dbutil.QuoteFn(backend)
	var timeRange offline.TimeRange
	query := fmt.Sprintf(`
		SELECT
			MIN(%s) AS %s,
			MAX(%s) AS %s
		FROM %s`, qt("unix_milli"), qt("min_unix_milli"), qt("unix_milli"), qt("max_unix_milli"), qt(tableName))

	if err := tx.GetContext(ctx, &timeRange, query); err != nil {
		return nil, errdefs.WithStack(err)
	}
	return &timeRange, nil
}
