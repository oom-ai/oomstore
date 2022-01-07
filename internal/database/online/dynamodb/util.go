package dynamodb

import (
	"time"

	"github.com/pkg/errors"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func serializeByTag(i interface{}, valueType types.ValueType) (out interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("failed to serailize by tag: %v", r)
		}
	}()

	switch valueType {
	case types.Time:
		return i.(time.Time).UnixMilli(), nil
	default:
		return i, nil
	}
}
