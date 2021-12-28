package sqlutil

import (
	"fmt"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const BatchSize = 10

func OnlineBatchTableName(revisionID int) string {
	return fmt.Sprintf("online_batch_%d", revisionID)
}

func OnlineStreamTableName(groupID int) string {
	return fmt.Sprintf("online_stream_%d", groupID)
}

func deserializeByTag(i interface{}, valueType types.ValueType, backend types.BackendType) (interface{}, error) {
	if i == nil {
		return nil, nil
	}

	switch valueType {
	case types.String:
		if backend == types.BackendMySQL {
			return string(i.([]byte)), nil
		}
		return i, nil
	case types.Bool:
		if backend == types.BackendMySQL || backend == types.BackendSQLite {
			if i == int64(1) {
				return true, nil
			} else if i == int64(0) {
				return false, nil
			} else {
				return nil, fmt.Errorf("invalid bool value: %s", i)
			}
		}
		return i, nil
	default:
		return i, nil
	}
}
