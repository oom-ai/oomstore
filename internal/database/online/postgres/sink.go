package postgres

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (s *DB) SinkFeatureValuesStream(ctx context.Context, stream <-chan *types.RawFeatureValueRecord, features []*types.Feature, revision *types.Revision) error {
	panic("not implemented")
}
