package kvutil

import (
	"fmt"
	"strconv"
	"time"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const serializeIntBase = 36

func SerializeByValueType(i interface{}, valueType types.ValueType) (s string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to serailize by tag: %v", r)
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
		return "", fmt.Errorf("unable to serialize %#v of type %T to string", i, i)
	}
}

func SerializeByValue(i interface{}) (string, error) {
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
		return SerializeByValue(s.UnixMilli())
	case bool:
		if s {
			return "1", nil
		} else {
			return "0", nil
		}

	default:
		return "", fmt.Errorf("unable to serialize %#v of type %T to string", i, i)
	}
}

func DeserializeByValueType(i interface{}, valueType types.ValueType) (interface{}, error) {
	if i == nil {
		return nil, nil
	}

	s, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("not a string or nil: %v", i)
	}

	switch valueType {
	case types.String:
		return s, nil

	case types.Int64:
		x, err := strconv.ParseInt(s, serializeIntBase, 64)
		return x, err

	case types.Float64:
		x, err := strconv.ParseFloat(s, 64)
		return x, err

	case types.Bool:
		if s == "1" {
			return true, nil
		} else if s == "0" {
			return false, nil
		} else {
			return nil, fmt.Errorf("invalid bool value: %s", s)
		}
	case types.Time:
		x, err := strconv.ParseInt(s, serializeIntBase, 64)
		return time.UnixMilli(x), err

	case types.Bytes:
		return []byte(s), nil
	default:
		return "", fmt.Errorf("unsupported value type: %s", valueType)
	}
}
