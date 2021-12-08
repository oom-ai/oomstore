package dynamodb_test

import (
	"context"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/internal/database/online/dynamodb"
	"github.com/ethhte88/oomstore/internal/database/online/test_impl"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

// The test depends on running amazon/dynamodb-local via
// docker run -d -p 8000:8000 amazon/dynamodb-local
// TODO: use a localstack dynamodb instance instead, check out https://github.com/elgohr/go-localstack
// gnomock does not support dynamodb yet, see https://github.com/orlangure/gnomock/issues/53
func prepareStore() (context.Context, online.Store) {
	db, err := dynamodb.Open(&types.DynamoDBOpt{
		Region:      "us-east-1",
		EndpointURL: "http://localhost:8000",
		// Test against local dynamodb instance doesn't rely upon credentials
		AccessKeyID:     "dummy",
		SecretAccessKey: "dummy",
		SessionToken:    "dummy",
		Source:          "dummy",
	})
	if err != nil {
		panic(err)
	}
	return context.Background(), db
}

func TestOpen(t *testing.T) {
	test_impl.TestOpen(t, prepareStore)
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore)
}
