package sqlutil

import (
	"context"
	"fmt"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"
)

func Export(ctx context.Context, db *sqlx.DB, opt offline.ExportOpt, backendType types.BackendType) (<-chan types.ExportRecord, <-chan error) {
	fields := append([]string{opt.EntityName}, opt.Features.Names()...)
	var fieldStr string
	var tableName string
	switch backendType {
	case types.POSTGRES, types.SNOWFLAKE:
		fieldStr = dbutil.Quote(`"`, fields...)
		tableName = dbutil.Quote(`"`, opt.DataTable)
	case types.MYSQL:
		fieldStr = dbutil.Quote("`", fields...)
		tableName = dbutil.Quote("`", opt.DataTable)
	}
	query := fmt.Sprintf("SELECT %s FROM %s", fieldStr, tableName)
	if opt.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *opt.Limit)
	}
	stream := make(chan types.ExportRecord)
	errs := make(chan error, 1) // at most 1 error

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
				if f.ValueType != types.STRING || record[i+1] == nil {
					continue
				}
				record[i+1] = cast.ToString(record[i+1])
			}
			stream <- record
		}
	}()

	return stream, errs
}
