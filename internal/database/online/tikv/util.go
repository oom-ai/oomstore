package tikv

import (
	"strings"

	"github.com/oom-ai/oomstore/internal/database/online/kvutil"
)

const keyDelimiter = ':'

func getKeyOfBatchFeature(revisionID, entityKey, featureID string) []byte {
	return []byte(kvutil.KeyPrefixForBatchFeature + getKey(revisionID, entityKey, featureID))
}

func getKeyOfStreamFeature(groupID, entityKey, featureID string) []byte {
	return []byte(kvutil.KeyPrefixForStreamFeature + getKey(groupID, entityKey, featureID))
}

func getKey(revisionID, entityKey, featureID string) string {
	return strings.Join([]string{revisionID, entityKey, featureID}, string(keyDelimiter))
}
