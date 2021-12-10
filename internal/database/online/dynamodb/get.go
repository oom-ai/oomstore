package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/internal/database/online/sqlutil"
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

	if result.Item == nil {
		return nil, nil
	}

	rowMap := make(dbutil.RowMap)
	var value interface{}
	for _, feature := range opt.FeatureList {
		attributeValue, ok := result.Item[feature.Name]
		if !ok {
			return nil, fmt.Errorf("could not find feature %s in table %s", feature.Name, tableName)
		}
		if err = attributevalue.Unmarshal(attributeValue, &value); err != nil {
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

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	res := make(map[string]dbutil.RowMap)
	for _, entityKey := range opt.EntityKeys {
		rowMap, err := db.Get(ctx, online.GetOpt{
			Entity:      opt.Entity,
			RevisionID:  opt.RevisionID,
			EntityKey:   entityKey,
			FeatureList: opt.FeatureList,
		})
		if err != nil {
			return res, err
		}
		if len(rowMap) > 0 {
			res[entityKey] = rowMap
		}
	}
	return res, nil
}
