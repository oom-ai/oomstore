package postgres

import (
	"fmt"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func DeserializeByTag(i interface{}, typeTag string) (interface{}, error) {
	if i == nil {
		return nil, nil
	}

	switch typeTag {
	case types.INT8:
		x, ok := i.(int64)
		if !ok {
			return nil, fmt.Errorf("%#v of type %T is not an int64", i, i)
		}
		return int8(x), nil
	case types.INT16:
		x, ok := i.(int64)
		if !ok {
			return nil, fmt.Errorf("%#v of type %T is not an int64", i, i)
		}
		return int16(x), nil
	case types.INT32:
		x, ok := i.(int64)
		if !ok {
			return nil, fmt.Errorf("%#v of type %T is not an int64", i, i)
		}
		return int32(x), nil

	case types.FLOAT32:
		x, ok := i.(float64)
		if !ok {
			return nil, fmt.Errorf("%#v of type %T is not an float64", i, i)
		}
		return float32(x), nil
	default:
		return i, nil
	}
}
