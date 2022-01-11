package dbutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func DeserializeByValueType(i interface{}, valueType types.ValueType, backend types.BackendType) (interface{}, error) {
	if i == nil {
		return nil, nil
	}
	var deserializer func(i interface{}, valueType types.ValueType) (interface{}, error)
	switch backend {
	case types.BackendCassandra:
		deserializer = cassandraDeserializer
	case types.BackendDynamoDB:
		deserializer = dynamoDeserializer
	case types.BackendSnowflake:
		deserializer = snowflakeDeserializer
	case types.BackendRedis, types.BackendTiKV:
		deserializer = kvDeserializer
	case types.BackendMySQL, types.BackendPostgres, types.BackendSQLite,
		types.BackendBigQuery, types.BackendRedshift:
		deserializer = rdbDeserializer
	default:
		panic(fmt.Sprintf("unsupported backend type %s", backend))
	}

	value, err := deserializer(i, valueType)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	return value, nil
}

func rdbDeserializer(i interface{}, valueType types.ValueType) (interface{}, error) {
	switch valueType {
	case types.Bool:
		s := strings.ToLower(cast.ToString(i))
		if s == "1" || s == "true" {
			return true, nil
		} else if s == "0" || s == "false" {
			return false, nil
		}
		return nil, errdefs.Errorf("invalid bool value %v", i)
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

func dynamoDeserializer(i interface{}, valueType types.ValueType) (interface{}, error) {
	switch valueType {
	case types.Int64:
		v, ok := i.(float64)
		if !ok {
			return "", errdefs.Errorf("not float64 %v", i)
		}
		return int64(v), nil
	case types.Time:
		v, ok := i.(float64)
		if !ok {
			return "", errdefs.Errorf("not float64 %v", i)
		}
		return time.UnixMilli(int64(v)), nil
	default:
		return i, nil
	}
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
		return nil, errdefs.Errorf("not a string or nil: %v", i)
	}

	switch valueType {
	case types.String:
		return s, nil

	case types.Int64:
		x, err := strconv.ParseInt(s, 10, 64)
		return x, errdefs.WithStack(err)

	case types.Float64:
		x, err := strconv.ParseFloat(s, 64)
		return x, errdefs.WithStack(err)

	case types.Bool:
		if s == "1" {
			return true, nil
		} else if s == "0" {
			return false, nil
		} else {
			return nil, errdefs.Errorf("invalid bool value: %s", s)
		}
	case types.Time:
		x, err := strconv.ParseInt(s, 10, 64)
		return time.UnixMilli(x), err

	case types.Bytes:
		return []byte(s), nil
	default:
		return "", errdefs.Errorf("unsupported value type: %s", valueType)
	}
}

func kvDeserializer(i interface{}, valueType types.ValueType) (interface{}, error) {
	if i == nil {
		return nil, nil
	}

	s, ok := i.(string)
	if !ok {
		return nil, errdefs.Errorf("not a string or nil: %v", i)
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
			return nil, errdefs.Errorf("invalid bool value: %s", s)
		}
	case types.Time:
		x, err := strconv.ParseInt(s, serializeIntBase, 64)
		return time.UnixMilli(x), err

	case types.Bytes:
		return []byte(s), nil
	default:
		return "", errdefs.Errorf("unsupported value type: %s", valueType)
	}
}
