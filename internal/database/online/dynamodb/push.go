package dynamodb

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/online"
)

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	var (
		tableName = dbutil.OnlineStreamTableName(opt.GroupID)
		item      = make(map[string]types.AttributeValue)
	)

	entityKeyValue, err := attributevalue.Marshal(opt.EntityKey)
	if err != nil {
		return errdefs.WithStack(err)
	}
	item[opt.EntityName] = entityKeyValue

	for i, feature := range opt.Features {
		value, err2 := dbutil.SerializeByValueType(opt.FeatureValues[i], feature.ValueType, Backend)
		if err2 != nil {
			return err2
		}
		attributeValue, err2 := attributevalue.Marshal(value)
		if err2 != nil {
			return errdefs.WithStack(err2)
		}
		item[feature.Name] = attributeValue
	}

	_, err = db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})
	return errdefs.WithStack(err)
}
