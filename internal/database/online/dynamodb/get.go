package dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	oomTypes "github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	BatchGetItemCapacity = 100
)

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	if err := opt.Validate(); err != nil {
		return nil, err
	}

	var tableName string
	if opt.Group.Category == oomTypes.CategoryBatch {
		tableName = dbutil.OnlineBatchTableName(*opt.RevisionID)
	} else {
		tableName = dbutil.OnlineStreamTableName(opt.Group.ID)
	}

	entityKeyValue, err := attributevalue.Marshal(opt.EntityKey)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	result, err := db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			opt.Group.Entity.Name: entityKeyValue,
		},
	})
	if err != nil {
		if apiErr := new(types.ResourceNotFoundException); errors.As(err, &apiErr) {
			return make(dbutil.RowMap), nil
		}
		return nil, errdefs.WithStack(err)
	}
	return deserializeFeatureValues(opt.Features, result.Item)
}

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	if err := opt.Validate(); err != nil {
		return nil, err
	}

	var tableName string
	if opt.Group.Category == oomTypes.CategoryBatch {
		tableName = dbutil.OnlineBatchTableName(*opt.RevisionID)
	} else {
		tableName = dbutil.OnlineStreamTableName(opt.Group.ID)
	}

	entityName := opt.Group.Entity.Name
	res := make(map[string]dbutil.RowMap)
	keys := make([]map[string]types.AttributeValue, 0, BatchGetItemCapacity)
	for _, entityKey := range opt.EntityKeys {
		entityKeyValue, err := attributevalue.Marshal(entityKey)
		if err != nil {
			return nil, errdefs.WithStack(err)
		}
		keys = append(keys, map[string]types.AttributeValue{
			entityName: entityKeyValue,
		})
		if len(keys) == BatchGetItemCapacity {
			if err = batchGetItem(ctx, db, keys, tableName, entityName, opt.Features, res); err != nil {
				return nil, err
			}
			keys = make([]map[string]types.AttributeValue, 0, BatchGetItemCapacity)
		}
	}
	if err := batchGetItem(ctx, db, keys, tableName, entityName, opt.Features, res); err != nil {
		return nil, err
	}
	return res, nil
}

func batchGetItem(ctx context.Context, db *DB, keys []map[string]types.AttributeValue, tableName, entityName string, features oomTypes.FeatureList, res map[string]dbutil.RowMap) error {
	result, err := db.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			tableName: {
				Keys: keys,
			},
		},
	})
	if err != nil {
		if apiErr := new(types.ResourceNotFoundException); !errors.As(err, &apiErr) {
			return errdefs.WithStack(err)
		}
	}

	if result == nil {
		return nil
	}

	for _, item := range result.Responses[tableName] {
		entityKeyValue, ok := item[entityName]
		if !ok {
			return errdefs.Errorf("could not find entity key column %s in table %s", entityName, tableName)
		}
		var entityKey string
		if err = attributevalue.Unmarshal(entityKeyValue, &entityKey); err != nil {
			return errdefs.WithStack(err)
		}
		rowMap, err := deserializeFeatureValues(features, item)
		if err != nil {
			return err
		}
		res[entityKey] = rowMap
	}
	return nil
}

func deserializeFeatureValues(features oomTypes.FeatureList, item map[string]types.AttributeValue) (dbutil.RowMap, error) {
	if item == nil {
		return make(dbutil.RowMap), nil
	}
	rowMap := make(dbutil.RowMap)
	var value interface{}
	for _, feature := range features {
		attributeValue, ok := item[feature.Name]
		if !ok {
			return nil, errdefs.Errorf("could not find feature %s", feature.Name)
		}
		if err := attributevalue.Unmarshal(attributeValue, &value); err != nil {
			return nil, errdefs.WithStack(err)
		}
		deserializedValue, err := dbutil.DeserializeByValueType(value, feature.ValueType, oomTypes.BackendDynamoDB)
		if err != nil {
			return nil, err
		}
		rowMap[feature.FullName()] = deserializedValue
	}
	return rowMap, nil
}
