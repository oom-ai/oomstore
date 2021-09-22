package cmd

import (
	"fmt"
)

var validCategories = []string{
	"batch",
	"stream",
}

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

func validateCategory(value string) error {
	return validateEnum("category", value, validCategories)
}

func validateStatus(value string) error {
	return validateEnum("status", value, validStatus)
}
