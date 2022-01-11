package tikv

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"

	"github.com/pkg/errors"

	"github.com/oom-ai/oomstore/internal/database/online/kvutil"
)

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	serializedRevisionID, err := dbutil.SerializeByValue(revisionID, Backend)
	if err != nil {
		return err
	}
	startKey := append([]byte(kvutil.KeyPrefixForBatchFeature+serializedRevisionID), byte(keyDelimiter))
	endKey := append([]byte(kvutil.KeyPrefixForBatchFeature+serializedRevisionID), byte(keyDelimiter+1))

	err = db.DeleteRange(ctx, startKey, endKey)
	return errors.WithStack(err)
}
