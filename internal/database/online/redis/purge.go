package redis

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"

	"github.com/oom-ai/oomstore/pkg/errdefs"
)

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	prefix, err := dbutil.SerializeByValue(revisionID, Backend)
	if err != nil {
		return nil
	}
	pattern := "b" + prefix + ":*"

	var cursor uint64
	var keys []string
	for {
		keys, cursor, err = db.Scan(ctx, cursor, pattern, PipelineBatchSize).Result()
		if err != nil {
			return errdefs.WithStack(err)
		}

		if len(keys) > 0 {
			if _, err = db.Del(ctx, keys...).Result(); err != nil {
				return errdefs.WithStack(err)
			}
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}
