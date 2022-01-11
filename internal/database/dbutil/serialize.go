package dbutil

import (
	"fmt"
	"strconv"
	"time"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	serializeIntBase = 36
)

func SerializeByValueType(i interface{}, valueType types.ValueType, backend types.BackendType) (interface{}, error) {
	if i == nil {
		return nil, nil
	}
	var serializer func(i interface{}, valueType types.ValueType) (interface{}, error)
	switch backend {
	case types.BackendRedis, types.BackendTiKV:
		serializer = kvSerializerByValueType
	case types.BackendDynamoDB:
		serializer = dynamoSerializerByValueType
	default:
		panic(fmt.Sprintf("unsupported backend type %s", backend))
	}

	value, err := serializer(i, valueType)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	return value, nil
}

func kvSerializerByValueType(i interface{}, valueType types.ValueType) (s interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errdefs.Errorf("failed to serailize by tag: %v", r)
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
		return "", errdefs.Errorf("unable to serialize %#v of type %T to string", i, i)
	}
}

func dynamoSerializerByValueType(i interface{}, valueType types.ValueType) (out interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errdefs.Errorf("failed to serailize by value type: %v", r)
		}
	}()

	switch valueType {
	case types.Time:
		return i.(time.Time).UnixMilli(), nil
	default:
		return i, nil
	}
}

func SerializeByValue(i interface{}, backend types.BackendType) (string, error) {
	if i == nil {
		return "", nil
	}
	var serializer func(i interface{}) (string, error)
	switch backend {
	case types.BackendRedis, types.BackendTiKV:
		serializer = kvSerializerByValue
	default:
		panic(fmt.Sprintf("unsupported backend type %s", backend))
	}

	value, err := serializer(i)
	if err != nil {
		return "", errdefs.WithStack(err)
	}
	return value, nil
}

func kvSerializerByValue(i interface{}) (string, error) {
	switch s := i.(type) {
	case string:
		return s, nil
	case []byte:
		return string(s), nil

	case int:
		return strconv.FormatInt(int64(s), serializeIntBase), nil
	case int64:
		return strconv.FormatInt(int64(s), serializeIntBase), nil
	case int32:
		return strconv.FormatInt(int64(s), serializeIntBase), nil
	case int16:
		return strconv.FormatInt(int64(s), serializeIntBase), nil
	case int8:
		return strconv.FormatInt(int64(s), serializeIntBase), nil

	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil

	case uint:
		return strconv.FormatUint(uint64(s), serializeIntBase), nil
	case uint64:
		return strconv.FormatUint(uint64(s), serializeIntBase), nil
	case uint32:
		return strconv.FormatUint(uint64(s), serializeIntBase), nil
	case uint16:
		return strconv.FormatUint(uint64(s), serializeIntBase), nil
	case uint8:
		return strconv.FormatUint(uint64(s), serializeIntBase), nil

	case time.Time:
		return kvSerializerByValue(s.UnixMilli())
	case bool:
		if s {
			return "1", nil
		} else {
			return "0", nil
		}

	default:
		return "", errdefs.Errorf("unable to serialize %#v of type %T to string", i, i)
	}
}
