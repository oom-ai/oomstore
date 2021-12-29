package tikv

import "strings"

const keyDelimiter = ":"

func getKey(revisionID string, entityKey string, featureID string) []byte {
	return []byte(strings.Join([]string{revisionID, entityKey, featureID}, keyDelimiter))
}
