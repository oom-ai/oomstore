package dbutil

import (
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/pkg/errors"
)

func DeserializeByValueType(i interface{}, valueType types.ValueType, backend types.BackendType) (interface{}, error) {
	if i == nil {
		return nil, nil
	}
	var deserializer func(i interface{}, valueType types.ValueType) (interface{}, error)
	switch backend {
	case types.BackendCassandra:
		deserializer = cassandraDeserializer
	case types.BackendSnowflake:
		deserializer = snowflakeDeserializer
	default:
		deserializer = defaultDeserializer
	}

	value, err := deserializer(i, valueType)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return value, nil
}

func defaultDeserializer(i interface{}, valueType types.ValueType) (interface{}, error) {
	switch valueType {
	case types.Bool:
		s := strings.ToLower(cast.ToString(i))
		if s == "1" || s == "true" {
			return true, nil
		} else if s == "0" || s == "false" {
			return false, nil
		}
		return nil, errors.Errorf("invalid bool value %v", i)
	case types.String:
		return cast.ToString(i), nil
	default:
		return i, nil
	}
}

func cassandraDeserializer(i interface{}, valueType types.ValueType) (interface{}, error) {
	switch i.(type) {
	case string:
		if i == "" {
			return nil, nil
		}
	}
	return i, nil
}

func snowflakeDeserializer(i interface{}, valueType types.ValueType) (interface{}, error) {
	if i == nil {
		return nil, nil
	}
	if valueType == types.Bool {
		return i, nil
	}

	s, ok := i.(string)
	if !ok {
		return nil, errors.Errorf("not a string or nil: %v", i)
	}

	switch valueType {
	case types.String:
		return s, nil

	case types.Int64:
		x, err := strconv.ParseInt(s, 10, 64)
		return x, errors.WithStack(err)

	case types.Float64:
		x, err := strconv.ParseFloat(s, 64)
		return x, errors.WithStack(err)

	case types.Bool:
		if s == "1" {
			return true, nil
		} else if s == "0" {
			return false, nil
		} else {
			return nil, errors.Errorf("invalid bool value: %s", s)
		}
	case types.Time:
		x, err := strconv.ParseInt(s, 10, 64)
		return time.UnixMilli(x), err

	case types.Bytes:
		return []byte(s), nil
	default:
		return "", errors.Errorf("unsupported value type: %s", valueType)
	}
}
