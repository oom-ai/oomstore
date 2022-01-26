package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/errdefs"
)

func (db *DB) CreateTable(ctx context.Context, opt online.CreateTableOpt) error {
	_, err := db.Client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(opt.TableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(opt.EntityName),
				KeyType:       types.KeyTypeHash,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(opt.EntityName),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	if err != nil {
		return errdefs.WithStack(err)
	}
	return nil
}
