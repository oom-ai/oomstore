package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) GetFeatureValuesStream(ctx context.Context, opt types.GetFeatureValuesStreamOpt) (<-chan *types.RawFeatureValueRecord, error) {
	query := fmt.Sprintf("select %s from %s", strings.Join(opt.FeatureNames, ","), opt.FeatureNames)
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
