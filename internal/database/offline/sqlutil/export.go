package sqlutil

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type QueryExportResults func(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.ExportOpt, query string, args []interface{}) (<-chan types.ExportRecord, <-chan error)

type DoExportOpt struct {
	offline.ExportOpt
	QueryResults QueryExportResults
}

func Export(ctx context.Context, db *sqlx.DB, opt offline.ExportOpt, backend types.BackendType) (<-chan types.ExportRecord, <-chan error) {
	dbOpt := dbutil.DBOpt{
		Backend: backend,
		SqlxDB:  db,
	}
	doJoinOpt := DoExportOpt{
		ExportOpt:    opt,
		QueryResults: sqlxQueryExportResults,
	}
	return DoExport(ctx, dbOpt, doJoinOpt)
}

func DoExport(ctx context.Context, dbOpt dbutil.DBOpt, opt DoExportOpt) (<-chan types.ExportRecord, <-chan error) {
	var (
		emptyStream = make(chan types.ExportRecord)
		errs        = make(chan error, 1) // at most 1 error
	)
	defer close(emptyStream)
	defer close(errs)

	query, args, err := buildExportQuery(dbOpt, opt.ExportOpt)
	if err != nil {
		errs <- errdefs.WithStack(err)
		return emptyStream, errs
	}

	return opt.QueryResults(ctx, dbOpt, opt.ExportOpt, query, args)
}

func sqlxQueryExportResults(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.ExportOpt, query string, args []interface{}) (<-chan types.ExportRecord, <-chan error) {
	stream := make(chan types.ExportRecord)
	errs := make(chan error, 1) // at most 1 error

	go func() {
		defer close(stream)
		defer close(errs)
		stmt, err := dbOpt.SqlxDB.Preparex(dbOpt.SqlxDB.Rebind(query))
		if err != nil {
			errs <- errdefs.WithStack(err)
			return
		}
		defer stmt.Close()
		rows, err := stmt.Queryx(args...)
		if err != nil {
			errs <- errdefs.WithStack(err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			record, err := rows.SliceScan()
			if err != nil {
				errs <- errdefs.Errorf("failed at rows.SliceScan, err=%v", err)
				return
			}
			record[0] = cast.ToString(record[0])
			for i, f := range opt.Features {
				if record[i+1] == nil {
					continue
				}
				deserializedValue, err := dbutil.DeserializeByValueType(record[i+1], f.ValueType, dbOpt.Backend)
				if err != nil {
					errs <- err
					return
				}
				record[i+1] = deserializedValue
			}
			stream <- record
		}
	}()

	return stream, errs
}

func buildExportQuery(dbOpt dbutil.DBOpt, opt offline.ExportOpt) (string, []interface{}, error) {
	if opt.CdcTable == nil && opt.UnixMilli == nil {
		return buildExportBatchQuery(dbOpt, opt), nil, nil
	}
	if opt.CdcTable != nil && opt.UnixMilli != nil {
		return buildExportStreamQuery(dbOpt, opt)
	}
	return "", nil, fmt.Errorf("invalid option %+v", opt)
}

func buildExportBatchQuery(dbOpt dbutil.DBOpt, opt offline.ExportOpt) string {
	fields := append([]string{opt.EntityName}, opt.Features.Names()...)
	qt := dbutil.QuoteFn(dbOpt.Backend)

	tableName := qt(opt.SnapshotTable)
	if dbOpt.Backend == types.BackendBigQuery {
		tableName = fmt.Sprintf("%s.%s", *dbOpt.DatasetID, tableName)
	}
	query := fmt.Sprintf("SELECT %s FROM %s", qt(fields...), tableName)
	if opt.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *opt.Limit)
	}
	return query
}

func buildExportStreamQuery(dbOpt dbutil.DBOpt, opt offline.ExportOpt) (string, []interface{}, error) {
	query, err := buildAggregateQuery(aggregateQueryParams{
		EntityName:            opt.EntityName,
		UnixMilli:             "unix_milli",
		FeatureNames:          opt.Features.Names(),
		PrevSnapshotTableName: opt.SnapshotTable,
		CurrCdcTableName:      *opt.CdcTable,
		Backend:               dbOpt.Backend,
		DatasetID:             dbOpt.DatasetID,
	})
	if err != nil {
		return "", nil, errdefs.WithStack(err)
	}
	if opt.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *opt.Limit)
	}
	return query, []interface{}{*opt.UnixMilli, *opt.UnixMilli, *opt.UnixMilli, *opt.UnixMilli}, nil
}
