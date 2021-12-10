package sqlutil

import (
	"fmt"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

const BatchSize = 10

func OnlineTableName(revisionID int) string {
	return fmt.Sprintf("online_%d", revisionID)
}

func deserializeByTag(i interface{}, typeTag string, backend types.BackendType) (interface{}, error) {
	if i == nil {
		return nil, nil
	}

	switch typeTag {
	case types.STRING:
		if backend == types.MYSQL {
			return string(i.([]byte)), nil
		}
		return i, nil
	case types.BOOL:
		if backend == types.MYSQL {
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
