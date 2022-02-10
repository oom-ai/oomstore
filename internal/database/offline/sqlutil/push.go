package sqlutil

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/errdefs"
)

func Push(ctx context.Context, dbOpt dbutil.DBOpt, pushOpt offline.PushOpt) error {
	tableName := dbutil.OfflineStreamCdcTableName(pushOpt.GroupID, pushOpt.Revision)

	columns := make([]string, 0, len(pushOpt.FeatureNames)+2)
	columns = append(columns, pushOpt.EntityName, "unix_milli")
	columns = append(columns, pushOpt.FeatureNames...)

	rows := make([]interface{}, 0, len(pushOpt.Records))
	for _, record := range pushOpt.Records {
		rows = append(rows, record.ToRow())
	}

	err := dbutil.InsertRecordsToTable(ctx, dbOpt, tableName, rows, columns)
	if err != nil {
		tableNotFound, notFoundErr := dbutil.IsTableNotFoundError(err, dbOpt.Backend)
		if notFoundErr != nil {
			return notFoundErr
		}
		if tableNotFound {
			return errdefs.NotFound(errdefs.WithStack(err))
		}
	}
	return errdefs.WithStack(err)
}
