package util

import (
	"strings"

	"github.com/oom-ai/oomstore/pkg/errdefs"
)

const SepFullFeatureName = "."

func ComposeFullFeatureName(groupName string, featureName string) string {
	return groupName + SepFullFeatureName + featureName
}

func SplitFullFeatureName(fullname string) (string, string, error) {
	parts := strings.SplitN(fullname, SepFullFeatureName, 2)
	if len(parts) != 2 {
		return "", "", errdefs.Errorf("invalid full feature name: '%s'", fullname)
	}
	return parts[0], parts[1], nil
}

func ValidateFullFeatureNames(fullnames ...string) error {
	for _, fullname := range fullnames {
		nameSlice := strings.Split(fullname, SepFullFeatureName)
		if len(nameSlice) != 2 {
			return errdefs.Errorf("invalid full feature name: '%s'", fullname)
		}
	}
	return nil
}
