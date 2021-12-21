package runtime_dynamodb

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/oom-ai/oomstore/internal/database/online/dynamodb"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/localstack"
)

var DynamoDBDbOpt types.DynamoDBOpt

func init() {
	dynamodbContainer, err := gnomock.Start(
		localstack.Preset(
			localstack.WithServices(localstack.DynamoDB),
		),
		gnomock.WithUseLocalImagesFirst(),
	)
	if err != nil {
		panic(err)
	}

	DynamoDBDbOpt = types.DynamoDBOpt{
		Region:          "us-east-1",
		EndpointURL:     fmt.Sprintf("http://%s/", dynamodbContainer.Address(localstack.APIPort)),
		AccessKeyID:     "test",
		SecretAccessKey: "test",
		SessionToken:    "test",
		Source:          "test",
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
		<-c

		_ = gnomock.Stop(dynamodbContainer)
	}()
}

func PrepareDB(t *testing.T) (context.Context, *dynamodb.DB) {
	db, err := dynamodb.Open(&DynamoDBDbOpt)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Drop all existing tables so that it doesn't interfere with tests that come after
	output, err := db.Client.ListTables(ctx, &awsDynamodb.ListTablesInput{})
	if err != nil {
		t.Fatal(err)
	}
	for _, tableName := range output.TableNames {
		if _, err := db.Client.DeleteTable(ctx, &awsDynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		}); err != nil {
			t.Fatal(err)
		}
	}

	return ctx, db
}
