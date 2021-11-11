package redis

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

func (db *DB) Purge(ctx context.Context, revision *typesv2.Revision) error {
	prefix, err := SerializeByValue(revision.ID)
	if err != nil {
		return nil
	}
	pattern := prefix + ":*"

	var cursor uint64
	var keys []string
	for {
		keys, cursor, err = db.Scan(ctx, cursor, pattern, PipelineBatchSize).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if _, err = db.Del(ctx, keys...).Result(); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}
