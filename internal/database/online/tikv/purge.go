package tikv

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/online/redis"
)

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	startKey, err := redis.SerializeByValue(revisionID)
	if err != nil {
		return err
	}
	endKey, err := redis.SerializeByValue(revisionID + 1)
	if err != nil {
		return err
	}
	return db.DeleteRange(ctx, []byte(startKey), []byte(endKey))
}
