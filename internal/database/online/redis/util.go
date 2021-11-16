package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func SerializeByTag(i interface{}, typeTag string) (s string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	switch typeTag {
	case types.STRING:
		return i.(string), nil
	case types.INT64:
		return strconv.FormatInt(int64(i.(int64)), SeralizeIntBase), nil
	case types.FLOAT64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64), nil
	case types.BOOL:
		if i.(bool) {
			return "1", nil
		} else {
			return "0", nil
		}
	case types.TIME:
		return strconv.FormatInt(i.(time.Time).UnixMilli(), SeralizeIntBase), nil

	case types.BYTE_ARRAY:
		return string(i.([]byte)), nil
	default:
		return "", fmt.Errorf("unable to seralize %#v of type %T to string", i, i)
	}
}

func SerializeByValue(i interface{}) (string, error) {
	switch s := i.(type) {
	case string:
		return s, nil
	case []byte:
		return string(s), nil

	case int:
		return strconv.FormatInt(int64(s), SeralizeIntBase), nil
	case int64:
		return strconv.FormatInt(int64(s), SeralizeIntBase), nil
	case int32:
		return strconv.FormatInt(int64(s), SeralizeIntBase), nil
	case int16:
		return strconv.FormatInt(int64(s), SeralizeIntBase), nil
	case int8:
		return strconv.FormatInt(int64(s), SeralizeIntBase), nil

	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil

	case uint:
		return strconv.FormatUint(uint64(s), SeralizeIntBase), nil
	case uint64:
		return strconv.FormatUint(uint64(s), SeralizeIntBase), nil
	case uint32:
		return strconv.FormatUint(uint64(s), SeralizeIntBase), nil
	case uint16:
		return strconv.FormatUint(uint64(s), SeralizeIntBase), nil
	case uint8:
		return strconv.FormatUint(uint64(s), SeralizeIntBase), nil

	case time.Time:
		return SerializeByValue(s.UnixMilli())
	case bool:
		if s {
			return "1", nil
		} else {
			return "0", nil
		}

	default:
		return "", fmt.Errorf("unable to seralize %#v of type %T to string", i, i)
	}
}

func DeserializeByTag(i interface{}, typeTag string) (interface{}, error) {
	if i == nil {
		return nil, nil
	}

	s, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("not a string or nil: %v", i)
	}

	switch typeTag {
	case types.STRING:
		return s, nil

	case types.INT64:
		x, err := strconv.ParseInt(s, SeralizeIntBase, 64)
		return x, err

	case types.FLOAT64:
		x, err := strconv.ParseFloat(s, 64)
		return x, err

	case types.BOOL:
		if s == "1" {
			return true, nil
		} else if s == "0" {
			return false, nil
		} else {
			return nil, fmt.Errorf("invalid bool value: %s", s)
		}
	case types.TIME:
		x, err := strconv.ParseInt(s, SeralizeIntBase, 64)
		return time.UnixMilli(x), err

	case types.BYTE_ARRAY:
		return []byte(s), nil
	default:
		return "", fmt.Errorf("unsupported type tag: %s", typeTag)
	}
}

func SerializeRedisKey(revisionId int32, entityKey interface{}) (string, error) {
	prefix, err := SerializeByValue(revisionId)
	if err != nil {
		return "", err
	}

	suffix, err := SerializeByValue(entityKey)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", prefix, suffix), nil
}
