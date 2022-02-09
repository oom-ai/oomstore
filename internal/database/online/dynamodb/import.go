package dynamodb

import (
	"context"
	"errors"

	"github.com/oom-ai/oomstore/internal/database/dbutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	oomTypes "github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	BatchWriteItemCapacity = 25
)

func (db *DB) Import(ctx context.Context, opt online.ImportOpt) error {
	// Step 0: clean up existing table for streaming feature
	var tableName string
	if opt.Group.Category == oomTypes.CategoryBatch {
		tableName = dbutil.OnlineBatchTableName(*opt.RevisionID)
	} else {
		tableName = dbutil.OnlineStreamTableName(opt.Group.ID)
		_, err := db.DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			if apiErr := new(types.ResourceNotFoundException); !errors.As(err, &apiErr) {
				return errdefs.WithStack(err)
			}
		}
	}
	// Step 1: create table
	if err := db.CreateTable(ctx, online.CreateTableOpt{
		TableName:  tableName,
		EntityName: opt.Group.Entity.Name,
	}); err != nil {
		return errdefs.WithStack(err)
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
	if err := batchWrite(ctx, db, tableName, items); err != nil {
		return err
	}

	if opt.ExportError != nil {
		select {
		case err := <-opt.ExportError:
			return err
		default:
			return nil
		}
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
		return nil, errdefs.WithStack(err)
	}
	item[opt.Group.Entity.Name] = entityKeyValue

	for i, feature := range opt.Features {
		value, err := dbutil.SerializeByValueType(record.ValueAt(i), feature.ValueType, Backend)
		if err != nil {
			return nil, err
		}
		attributeValue, err := attributevalue.Marshal(value)
		if err != nil {
			return nil, errdefs.WithStack(err)
		}
		item[feature.Name] = attributeValue
	}
	return item, nil
}
