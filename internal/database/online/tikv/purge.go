package tikv

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/online/kvutil"
)

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	startKey, err := kvutil.SerializeByValue(revisionID)
	if err != nil {
		return err
	}
	endKey, err := kvutil.SerializeByValue(revisionID + 1)
	if err != nil {
		return err
	}
	return db.DeleteRange(ctx, []byte(kvutil.KeyPrefixForBatchFeature+startKey), []byte(kvutil.KeyPrefixForBatchFeature+endKey))
}
