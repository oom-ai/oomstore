package dynamodb

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
)

func (db *DB) PrepareStreamTable(ctx context.Context, opt online.PrepareStreamTableOpt) error {
	// dynamodb has no "column", so we do nothing to "add column". see https://stackoverflow.com/a/25610645/16428442
	if opt.Feature != nil {
		return nil
	}

	tableName := sqlutil.OnlineStreamTableName(opt.GroupID)

	_, err := db.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(opt.Entity.Name),
				KeyType:       types.KeyTypeHash,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(opt.Entity.Name),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	return errors.WithStack(err)
}

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	var (
		tableName = sqlutil.OnlineStreamTableName(opt.GroupID)
		item      = make(map[string]types.AttributeValue)
	)

	entityKeyValue, err := attributevalue.Marshal(opt.EntityKey)
	if err != nil {
		return errors.WithStack(err)
	}
	item[opt.Entity.Name] = entityKeyValue

	for i, feature := range opt.Features {
		value, err := dbutil.SerializeByValueType(opt.FeatureValues[i], feature.ValueType, Backend)
		if err != nil {
			return err
		}
		attributevalue, err := attributevalue.Marshal(value)
		if err != nil {
			return errors.WithStack(err)
		}
		item[feature.Name] = attributevalue
	}

	_, err = db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})
	return errors.WithStack(err)
}
