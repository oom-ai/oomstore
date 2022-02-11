package cassandra

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"
)

const PUSH_QUERY = `
INSERT INTO {{ .TableName }} ( {{ .Fields }} )
VALUES ( {{ .InsertPlaceholders }} )
`

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	params := sqlutil.BuildPushQueryParams(opt, Backend)
	query, err := sqlutil.BuildPushQuery(params, PUSH_QUERY)
	if err != nil {
		return err
	}

	err = db.Query(query, params.InsertValues...).WithContext(ctx).Exec()
	return errdefs.WithStack(err)
}
