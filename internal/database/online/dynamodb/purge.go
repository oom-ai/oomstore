package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ethhte88/oomstore/internal/database/online/sqlutil"
)

func (db *DB) Purge(ctx context.Context, revisionID int) error {
	tableName := sqlutil.OnlineTableName(revisionID)
	_, err := db.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	return err
}
