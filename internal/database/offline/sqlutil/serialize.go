package sqlutil

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func DeserializeByValueType(i interface{}, valueType types.ValueType, backend types.BackendType) (interface{}, error) {
	if backend == types.BackendSnowflake {
		return deserializeByTagForSnowflake(i, valueType)
	}

	switch valueType {
	case types.Bool:
		s := cast.ToString(i)
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

// gosnowflake Scan always produce string when the destination is interface{}
// See https://github.com/snowflakedb/gosnowflake/issues/517
// As a work around, we cast the string to interface{} based on ValueType
// This method is mostly copied from redis.DeserializeByTag, except we use 10 rather than 36 as the base
// TODO: we should let the snowflake team fix the gosnowflake converter
func deserializeByTagForSnowflake(i interface{}, valueType types.ValueType) (interface{}, error) {
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
