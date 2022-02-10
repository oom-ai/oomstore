package util

import (
	"strings"

	"github.com/oom-ai/oomstore/pkg/errdefs"
)

const SepFullFeatureName = "."

func ComposeFullFeatureName(groupName string, featureName string) string {
	return groupName + SepFullFeatureName + featureName
}

func SplitFullFeatureName(fullName string) (string, string, error) {
	parts := strings.SplitN(fullName, SepFullFeatureName, 2)
	if len(parts) != 2 {
		return "", "", errdefs.Errorf("invalid full feature name: '%s'", fullName)
	}
	return parts[0], parts[1], nil
}

func ValidateFullFeatureNames(fullNames ...string) error {
	for _, fullName := range fullNames {
		nameSlice := strings.Split(fullName, SepFullFeatureName)
		if len(nameSlice) != 2 {
			return errdefs.Errorf("invalid full feature name: '%s'", fullName)
		}
	}
	return nil
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
