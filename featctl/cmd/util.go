package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/onestore-ai/onestore/pkg/onestore"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

var validStatus = []string{
	"enabled",
	"disabled",
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func validateEnum(name string, value string, possibleValues []string) error {
	if !contains(possibleValues, value) {
		return fmt.Errorf("invalid %s: '%s', should be one of %v", name, value, possibleValues)
	}
	return nil
}

func validateStatus(value string) error {
	return validateEnum("status", value, validStatus)
}

func mustOpenOneStore(ctx context.Context, opt types.OneStoreOpt) *onestore.OneStore {
	store, err := onestore.Open(ctx, oneStoreOpt)
	if err != nil {
		log.Fatalf("failed opening OneStore: %v", err)
	}
	return store
}
