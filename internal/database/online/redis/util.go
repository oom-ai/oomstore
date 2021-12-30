package redis

import (
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/online/kvutil"
)

func serializeRedisKey(revisionID int, entityKey interface{}) (string, error) {
	prefix, err := kvutil.SerializeByValue(revisionID)
	if err != nil {
		return "", err
	}

	suffix, err := kvutil.SerializeByValue(entityKey)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", prefix, suffix), nil
}
