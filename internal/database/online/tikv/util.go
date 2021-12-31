package tikv

import (
	"strings"

	"github.com/oom-ai/oomstore/internal/database/online/kvutil"
)

const keyDelimiter = ":"

func getKeyOfBatchFeature(revisionID string, entityKey string, featureID string) []byte {
	return []byte(kvutil.KeyPrefixForBatchFeature + getKey(revisionID, entityKey, featureID))
}

func getKeyOfStreamFeature(groupID string, entityKey string, featureID string) []byte {
	return []byte(kvutil.KeyPrefixForStreamFeature + getKey(groupID, entityKey, featureID))
}

func getKey(revisionID string, entityKey string, featureID string) string {
	return strings.Join([]string{revisionID, entityKey, featureID}, keyDelimiter)
}
