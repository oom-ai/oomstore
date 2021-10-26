package redis

import (
	"fmt"
	"strconv"
	"time"
)

// seralize feature values into compact string
func Seralize(i interface{}) (string, error) {
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
		return Seralize(s.UnixMilli())
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
