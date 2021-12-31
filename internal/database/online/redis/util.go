package redis

import (
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/online/kvutil"
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
	prefix, err := kvutil.SerializeByValue(prefixID)
	if err != nil {
		return "", err
	}

	suffix, err := kvutil.SerializeByValue(entityKey)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", prefix, suffix), nil
}
