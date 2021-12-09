package dynamodb

import (
	"fmt"
	"time"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func serializeByTag(i interface{}, typeTag string) (out interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to serailize by tag: %v", r)
		}
	}()

	switch typeTag {
	case types.TIME:
		return i.(time.Time).UnixMilli(), nil
	default:
		return i, nil
	}
}

func deserializeByTag(i interface{}, typeTag string) (interface{}, error) {
	if i == nil {
		return nil, nil
	}

	switch typeTag {
	case types.INT64:
		v, ok := i.(float64)
		if !ok {
			return "", fmt.Errorf("not float64 %v", i)
		}
		return int64(v), nil
	case types.TIME:
		v, ok := i.(float64)
		if !ok {
			return "", fmt.Errorf("not float64 %v", i)
		}
		return time.UnixMilli(int64(v)), nil
	default:
		return i, nil
	}
}
