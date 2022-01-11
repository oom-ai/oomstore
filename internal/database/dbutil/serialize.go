package dbutil

import (
	"fmt"
	"strconv"
	"time"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/pkg/errors"
)

const (
	serializeIntBase = 36
)

func SerializeByValueType(i interface{}, valueType types.ValueType, backend types.BackendType) (string, error) {
	if i == nil {
		return "", nil
	}
	var serializer func(i interface{}, valueType types.ValueType) (string, error)
	switch backend {
	case types.BackendRedis, types.BackendTiKV:
		serializer = kvSerializer
	default:
		return "", fmt.Errorf("unsupported backend type %s", backend)
	}

	value, err := serializer(i, valueType)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return value, nil
}

func kvSerializer(i interface{}, valueType types.ValueType) (s string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("failed to serailize by tag: %v", r)
		}
	}()

	switch valueType {
	case types.String:
		return i.(string), nil
	case types.Int64:
		return strconv.FormatInt(int64(i.(int64)), serializeIntBase), nil
	case types.Float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64), nil
	case types.Bool:
		if i.(bool) {
			return "1", nil
		} else {
			return "0", nil
		}
	case types.Time:
		return strconv.FormatInt(i.(time.Time).UnixMilli(), serializeIntBase), nil

	case types.Bytes:
		return string(i.([]byte)), nil
	default:
		return "", errors.Errorf("unable to serialize %#v of type %T to string", i, i)
	}
}
