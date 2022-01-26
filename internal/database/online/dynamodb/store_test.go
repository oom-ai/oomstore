package dynamodb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/test_impl"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_dynamodb"
)

func prepareStore(t *testing.T) (context.Context, online.Store) {
	return runtime_dynamodb.PrepareDB(t)
}

func destroyStore(t *testing.T) func() {
	return func() {
		ctx, db := runtime_dynamodb.PrepareDB(t)
		defer db.Close()

		// Drop all existing tables so that it doesn't interfere with tests that come after
		output, err := db.Client.ListTables(ctx, &awsDynamodb.ListTablesInput{})
		if err != nil {
			panic(err)
		}
		for _, tableName := range output.TableNames {
			if _, err := db.Client.DeleteTable(ctx, &awsDynamodb.DeleteTableInput{
				TableName: aws.String(tableName),
			}); err != nil {
				panic(fmt.Sprintf("failed deleting table '%s': %v", tableName, err))
			}
		}
	}
}

func TestOpen(t *testing.T) {
	test_impl.TestOpen(t, prepareStore, destroyStore(t))
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore, destroyStore(t))
}

func TestGetExisted(t *testing.T) {
	test_impl.TestGetExisted(t, prepareStore, destroyStore(t))
}

func TestGetNotExistedEntityKey(t *testing.T) {
	test_impl.TestGetNotExistedEntityKey(t, prepareStore, destroyStore(t))
}

func TestMultiGet(t *testing.T) {
	test_impl.TestMultiGet(t, prepareStore, destroyStore(t))
}

func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
	test_impl.TestPurgeRemovesSpecifiedRevision(t, prepareStore, destroyStore(t))
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
	test_impl.TestPurgeNotRemovesOtherRevisions(t, prepareStore, destroyStore(t))
}

func TestPush(t *testing.T) {
	test_impl.TestPush(t, prepareStore, destroyStore(t))
}

func TestCreateTable(t *testing.T) {
	test_impl.TestCreateTable(t, prepareStore, destroyStore(t))
}
