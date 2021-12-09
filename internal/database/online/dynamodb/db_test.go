package dynamodb_test

import (
	"context"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/internal/database/online/test_impl"
	"github.com/ethhte88/oomstore/internal/database/test/runtime_dynamodb"
)

func prepareStore() (context.Context, online.Store) {
	return runtime_dynamodb.PrepareDB()
}

func TestOpen(t *testing.T) {
	test_impl.TestOpen(t, prepareStore)
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore)
}
