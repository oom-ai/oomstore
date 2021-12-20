package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"
	oomTypes "github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	BatchGetItemCapacity = 100
)

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	tableName := sqlutil.OnlineTableName(opt.RevisionID)
	entityKeyValue, err := attributevalue.Marshal(opt.EntityKey)
	if err != nil {
		return nil, err
	}
	result, err := db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			opt.Entity.Name: entityKeyValue,
		},
	})
	if err != nil {
		if apiErr := new(types.ResourceNotFoundException); errors.As(err, &apiErr) {
			return nil, nil
		}
		return nil, err
	}
	return deserializeFeatureValues(opt.FeatureList, result.Item)
}

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	res := make(map[string]dbutil.RowMap)
	tableName := sqlutil.OnlineTableName(opt.RevisionID)
	keys := make([]map[string]types.AttributeValue, 0, BatchGetItemCapacity)
	for _, entityKey := range opt.EntityKeys {
		entityKeyValue, err := attributevalue.Marshal(entityKey)
		if err != nil {
			return nil, err
		}
		keys = append(keys, map[string]types.AttributeValue{
			opt.Entity.Name: entityKeyValue,
		})
		if len(keys) == BatchGetItemCapacity {
			if err = batchGetItem(ctx, db, keys, tableName, opt.Entity.Name, opt.FeatureList, res); err != nil {
				return nil, err
			}
			keys = make([]map[string]types.AttributeValue, 0, BatchGetItemCapacity)
		}
	}
	if err := batchGetItem(ctx, db, keys, tableName, opt.Entity.Name, opt.FeatureList, res); err != nil {
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
			return err
		}
	}

	for _, item := range result.Responses[tableName] {
		entityKeyValue, ok := item[entityName]
		if !ok {
			return fmt.Errorf("could not find entity key column %s in table %s", entityName, tableName)
		}
		var entityKey string
		if err = attributevalue.Unmarshal(entityKeyValue, &entityKey); err != nil {
			return err
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
		return nil, nil
	}
	rowMap := make(dbutil.RowMap)
	var value interface{}
	for _, feature := range features {
		attributeValue, ok := item[feature.Name]
		if !ok {
			return nil, fmt.Errorf("could not find feature %s", feature.Name)
		}
		if err := attributevalue.Unmarshal(attributeValue, &value); err != nil {
			return nil, err
		}
		typedValue, err := deserializeByTag(value, feature.ValueType)
		if err != nil {
			return nil, err
		}
		rowMap[feature.Name] = typedValue
	}
	return rowMap, nil
}
