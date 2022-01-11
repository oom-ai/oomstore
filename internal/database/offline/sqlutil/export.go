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

func Export(ctx context.Context, db *sqlx.DB, opt offline.ExportOpt, backend types.BackendType) (<-chan types.ExportRecord, <-chan error) {
	var (
		stream = make(chan types.ExportRecord)
		errs   = make(chan error, 1) // at most 1 error
	)

	query, args, err := buildExportQuery(opt, backend)
	if err != nil {
		defer close(stream)
		defer close(errs)
		errs <- errdefs.WithStack(err)
		return stream, errs
	}

	go func() {
		defer close(stream)
		defer close(errs)
		stmt, err := db.Preparex(db.Rebind(query))
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
				deserializedValue, err := dbutil.DeserializeByValueType(record[i+1], f.ValueType, backend)
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

func buildExportQuery(opt offline.ExportOpt, backend types.BackendType) (string, []interface{}, error) {
	if opt.CdcTable == nil && opt.UnixMilli == nil {
		return buildExportBatchQuery(opt, backend), nil, nil
	}
	if opt.CdcTable != nil && opt.UnixMilli != nil {
		return buildExportStreamQuery(opt, backend)
	}
	return "", nil, fmt.Errorf("invalid option %+v", opt)
}

func buildExportBatchQuery(opt offline.ExportOpt, backend types.BackendType) string {
	fields := append([]string{opt.EntityName}, opt.Features.Names()...)
	qt := dbutil.QuoteFn(backend)

	query := fmt.Sprintf("SELECT %s FROM %s", qt(fields...), qt(opt.SnapshotTable))
	if opt.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *opt.Limit)
	}
	return query
}

func buildExportStreamQuery(opt offline.ExportOpt, backend types.BackendType) (string, []interface{}, error) {
	query, err := buildAggregateQuery(aggregateQueryParams{
		EntityName:            opt.EntityName,
		UnixMilli:             "unix_milli",
		FeatureNames:          opt.Features.Names(),
		PrevSnapshotTableName: opt.SnapshotTable,
		CurrCdcTableName:      *opt.CdcTable,
		Backend:               backend,
	})
	if err != nil {
		return "", nil, errdefs.WithStack(err)
	}
	if opt.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *opt.Limit)
	}
	return query, []interface{}{*opt.UnixMilli, *opt.UnixMilli, *opt.UnixMilli, *opt.UnixMilli}, nil
}
