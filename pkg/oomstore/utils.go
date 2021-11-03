package oomstore

import (
	"strconv"

	"github.com/spf13/cast"
)

func castToInt64(i interface{}) (int64, error) {
	val := cast.ToInt(i)
	var err error
	if bytes, ok := i.([]byte); ok {
		val, err = strconv.Atoi(string(bytes))
		if err != nil {
			return 0, err
		}
	}
	return int64(val), nil
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func Int32Ptr(i int32) *int32 {
	return &i
}
