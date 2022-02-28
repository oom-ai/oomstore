package oomstore

import (
	"context"
)

func (o *OomStore) DropTemporaryTables(ctx context.Context, tableNames []string) error {
	return o.offline.DropTemporaryTable(ctx, tableNames)
}

func (o *OomStore) GetTemporaryTables(ctx context.Context, unixMilli int64) ([]string, error) {
	return o.offline.GetTemporaryTables(ctx, unixMilli)
}
