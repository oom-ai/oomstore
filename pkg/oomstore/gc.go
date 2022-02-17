package oomstore

import (
	"context"
	"time"
)

func (o *OomStore) Gc(ctx context.Context) error {
	return o.offline.DropTemporaryTable(ctx, time.Now().Add(time.Hour*24*-1).UnixMilli())
}
