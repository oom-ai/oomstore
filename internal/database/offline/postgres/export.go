package postgres

import (
	"context"
	"fmt"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func (db *DB) Export(ctx context.Context, opt offline.ExportOpt) (<-chan types.ExportRecord, <-chan error) {
	fields := append([]string{opt.EntityName}, opt.FeatureNames...)
	query := fmt.Sprintf("select %s from %s", dbutil.Quote(`"`, fields...), opt.DataTable)
	if opt.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *opt.Limit)
	}
	stream := make(chan types.ExportRecord)
	errs := make(chan error, 1) // at most 1 error

	go func() {
		defer close(stream)
		defer close(errs)
		rows, err := db.QueryxContext(ctx, query)
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
			stream <- record
		}
	}()

	return stream, errs
}
