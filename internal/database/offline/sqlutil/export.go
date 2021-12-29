package sqlutil

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
)

func Export(ctx context.Context, db *sqlx.DB, opt offline.ExportOpt, backendType types.BackendType) (<-chan types.ExportRecord, <-chan error) {
	var (
		fields = append([]string{opt.EntityName}, opt.Features.Names()...)
		stream = make(chan types.ExportRecord)
		errs   = make(chan error, 1) // at most 1 error
	)
	qt, err := dbutil.QuoteFn(backendType)
	if err != nil {
		errs <- err
		return stream, errs
	}
	query := fmt.Sprintf("SELECT %s FROM %s", qt(fields...), qt(opt.SnapshotTable))
	if opt.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *opt.Limit)
	}

	go func() {
		defer close(stream)
		defer close(errs)
		stmt, err := db.Preparex(query)
		if err != nil {
			errs <- err
			return
		}
		defer stmt.Close()
		rows, err := stmt.Queryx()
		if err != nil {
			errs <- err
			return
		}
		defer rows.Close()
		for rows.Next() {
			record, err := rows.SliceScan()
			if err != nil {
				errs <- fmt.Errorf("failed at rows.SliceScan, err=%v", err)
				return
			}
			record[0] = cast.ToString(record[0])
			for i, f := range opt.Features {
				if record[i+1] == nil {
					continue
				}
				if backendType == types.BackendSnowflake {

					v, err := deserializeByTagForSnowflake(record[i+1], f.ValueType)
					if err != nil {
						errs <- fmt.Errorf("failed at deserializeByTag, err=%v", err)
						return
					}
					record[i+1] = v
				} else {
					if f.ValueType == types.String {
						record[i+1] = cast.ToString(record[i+1])
					}
				}
			}
			stream <- record
		}
	}()

	return stream, errs
}
