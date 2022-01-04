package sqlutil

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/errdefs"
)

func Push(ctx context.Context, dbOpt dbutil.DBOpt, pushOpt offline.PushOpt) error {
	tableName := dbutil.OfflineStreamCdcTableName(pushOpt.GroupID, pushOpt.Revision)
	err := dbutil.InsertRecordsToTable(ctx, dbOpt, tableName, pushOpt.FeatureValues, pushOpt.FeatureNames)
	if err != nil && dbutil.IsTableNotFoundError(err, dbOpt.Backend) {
		return errdefs.NotFound(err)
	}
	return err
}
