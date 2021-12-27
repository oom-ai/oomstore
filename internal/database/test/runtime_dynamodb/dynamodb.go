package runtime_dynamodb

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online/dynamodb"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func init() {
	if out, err := exec.Command("oomplay", "init", "dynamodb").CombinedOutput(); err != nil {
		panic(fmt.Sprintf("oomplay failed with error: %v, output: %s", err, out))
	}
}

func PrepareDB(t *testing.T) (context.Context, *dynamodb.DB) {
	db, err := dynamodb.Open(GetOpt())
	if err != nil {
		t.Fatal(err)
	}
	return context.Background(), db
}

func GetOpt() *types.DynamoDBOpt {
	return &types.DynamoDBOpt{
		Region:          ".",
		EndpointURL:     "http://localhost:24566",
		AccessKeyID:     ".",
		SecretAccessKey: ".",
		SessionToken:    ".",
		Source:          ".",
	}
}
