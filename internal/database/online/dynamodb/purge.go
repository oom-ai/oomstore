package dynamodb

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	tableName := dbutil.OnlineBatchTableName(revisionID)
	_, err := db.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	return errdefs.WithStack(err)
}
