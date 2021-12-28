package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const BackendType = types.BackendDynamoDB

var _ online.Store = &DB{}

type DB struct {
	*dynamodb.Client
}

func Open(opt *types.DynamoDBOpt) (*DB, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(opt.Region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: opt.EndpointURL}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     opt.AccessKeyID,
				SecretAccessKey: opt.SecretAccessKey,
				SessionToken:    opt.SessionToken,
				Source:          opt.Source,
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	return &DB{dynamodb.NewFromConfig(cfg)}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	_, err := db.ListTables(ctx, nil)
	return err
}

// dynamodb is serverless so Close won't do anything
func (db *DB) Close() error {
	return nil
}

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	panic("Implement me!")
}
