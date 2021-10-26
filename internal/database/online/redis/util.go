package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
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

	case types.INT8:
		return strconv.FormatInt(int64(i.(int8)), SeralizeIntBase), nil
	case types.INT16:
		return strconv.FormatInt(int64(i.(int16)), SeralizeIntBase), nil
	case types.INT32:
		return strconv.FormatInt(int64(i.(int32)), SeralizeIntBase), nil
	case types.INT64:
		return strconv.FormatInt(int64(i.(int64)), SeralizeIntBase), nil

	case types.FLOAT32:
		return strconv.FormatFloat(float64(i.(float32)), 'f', -1, 32), nil
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
		return "", fmt.Errorf("unable to seralize %#v of type %T", i, i)
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
