package dynamodb

import (
	"fmt"
	"time"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func serializeByTag(i interface{}, valueType types.ValueType) (out interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to serailize by tag: %v", r)
		}
	}()

	switch valueType {
	case types.Time:
		return i.(time.Time).UnixMilli(), nil
	default:
		return i, nil
	}
}

func deserializeByTag(i interface{}, valueType types.ValueType) (interface{}, error) {
	if i == nil {
		return nil, nil
	}

	switch valueType {
	case types.Int64:
		v, ok := i.(float64)
		if !ok {
			return "", fmt.Errorf("not float64 %v", i)
		}
		return int64(v), nil
	case types.Time:
		v, ok := i.(float64)
		if !ok {
			return "", fmt.Errorf("not float64 %v", i)
		}
		return time.UnixMilli(int64(v)), nil
	default:
		return i, nil
	}
}
