package sqlutil

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/errdefs"
)

func cdcTableName(groupID int, revision int64) string {
	return fmt.Sprintf("offline_stream_cdc_%d_%d", groupID, revision)
}

func Push(ctx context.Context, dbOpt dbutil.DBOpt, pushOpt offline.PushOpt) error {
	tableName := cdcTableName(pushOpt.GroupID, pushOpt.Revision)
	err := dbutil.InsertRecordsToTable(ctx, dbOpt, tableName, pushOpt.FeatureValues, pushOpt.FeatureNames, dbOpt.Backend)
	if err != nil && dbutil.IsTableNotFoundError(err, dbOpt.Backend) {
		return errdefs.NotFound(err)
	}
	return err
}
