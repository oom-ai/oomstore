package postgres

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) SinkFeatureValuesStream(ctx context.Context, stream <-chan []interface{}, features []*types.Feature, revision *types.Revision) error {
	panic("not implemented")
}
