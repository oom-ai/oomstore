package dynamodb

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
	oomTypes "github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/pkg/errors"
)

const (
	BatchWriteItemCapacity = 25
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	// Step 1: create table
	tableName := sqlutil.OnlineBatchTableName(opt.Revision.ID)
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
	if err != nil {
		return errors.WithStack(err)
	}

	// Step 2: import items to the table
	items := make([]types.WriteRequest, 0, BatchWriteItemCapacity)
	for record := range opt.ExportStream {
		item, err := buildItem(record, opt)
		if err != nil {
			return err
		}
		items = append(items, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: item,
			},
		})
		if len(items) == BatchWriteItemCapacity {
			if err = batchWrite(ctx, db, tableName, items); err != nil {
				return err
			}
			items = make([]types.WriteRequest, 0, BatchWriteItemCapacity)
		}
	}
	if err = batchWrite(ctx, db, tableName, items); err != nil {
		return err
	}

	if opt.ExportError != nil {
		return <-opt.ExportError
	}
	return nil
}

func batchWrite(ctx context.Context, db *DB, tableName string, items []types.WriteRequest) error {
	if len(items) == 0 {
		return nil
	}
	_, err := db.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: items,
		},
	})
	return err
}

func buildItem(record oomTypes.ExportRecord, opt online.ImportOpt) (map[string]types.AttributeValue, error) {
	item := make(map[string]types.AttributeValue)
	entityKeyValue, err := attributevalue.Marshal(record.EntityKey())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	item[opt.Entity.Name] = entityKeyValue

	for i, feature := range opt.Features {
		value, err := dbutil.SerializeByValueType(record.ValueAt(i), feature.ValueType, Backend)
		if err != nil {
			return nil, err
		}
		attributeValue, err := attributevalue.Marshal(value)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		item[feature.Name] = attributeValue
	}
	return item, nil
}
