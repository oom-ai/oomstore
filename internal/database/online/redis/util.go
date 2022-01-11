package redis

import (
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
)

func serializeRedisKeyForBatchFeature(revisionID int, entityKey interface{}) (string, error) {
	s, err := serializeRediskey(revisionID, entityKey)
	return "b" + s, err
}

func serializeRedisKeyForStreamFeature(groupID int, entityKey interface{}) (string, error) {
	s, err := serializeRediskey(groupID, entityKey)
	return "s" + s, err
}

func serializeRediskey(prefixID int, entityKey interface{}) (string, error) {
	prefix, err := dbutil.SerializeByValue(prefixID, Backend)
	if err != nil {
		return "", err
	}

	suffix, err := dbutil.SerializeByValue(entityKey, Backend)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", prefix, suffix), nil
}
