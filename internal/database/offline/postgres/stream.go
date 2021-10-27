package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/onestore/types"
)

func (db *DB) Export(ctx context.Context, opt offline.ExportOpt) (<-chan *types.RawFeatureValueRecord, error) {
	fields := append([]string{opt.EntityName}, opt.FeatureNames...)
	query := fmt.Sprintf("select %s from %s", strings.Join(fields, ","), opt.DataTable)
	if opt.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *opt.Limit)
	}

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	stream := make(chan *types.RawFeatureValueRecord)
	go func() {
		defer rows.Close()
		defer close(stream)
		for rows.Next() {
			record, err := rows.SliceScan()
			stream <- &types.RawFeatureValueRecord{
				Record: record,
				Error:  err,
			}
			if err != nil {
				return
			}
		}
	}()

	return stream, nil
}
