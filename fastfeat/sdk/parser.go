package sdk

import (
	"strconv"
	"strings"
)

const sep = ","

func parseRawFeatureResponse(rawResponse string) Value {
	if rawResponse == FASTFEAT_EMPTY_VALUE {
		return EmptyValue{}
	}

	elementType := getRawFeatureResponseElementType(rawResponse)

	switch elementType {
	case Int64Type:
		res := Int64ArrayValue{}
		for _, e := range strings.Split(rawResponse, sep) {
			int64Value, _ := strconv.ParseInt(e, 10, 0)
			res.Int64ArrayValue = append(res.Int64ArrayValue, Int64Value{int64Value})
		}

		if len(res.Int64ArrayValue) == 1 {
			return res.Int64ArrayValue[0]
		}
		return res
	case DoubleType:
		res := DoubleArrayValue{}
		for _, e := range strings.Split(rawResponse, sep) {
			doubleValue, _ := strconv.ParseFloat(e, 64)
			res.DoubleArrayValue = append(res.DoubleArrayValue, DoubleValue{doubleValue})
		}

		if len(res.DoubleArrayValue) == 1 {
			return res.DoubleArrayValue[0]
		}
		return res
	default:
		res := StringArrayValue{}
		for _, e := range strings.Split(rawResponse, sep) {
			res.StringArrayValue = append(res.StringArrayValue, StringValue{e})
		}

		if len(res.StringArrayValue) == 1 {
			return res.StringArrayValue[0]
		}
		return res
	}
}

func getRawFeatureResponseElementType(rawResponse string) FeatureElementTypeEnum {
	mp := make(map[FeatureElementTypeEnum]int)

	vs := strings.Split(rawResponse, sep)
	for _, v := range vs {
		if _, err := strconv.ParseInt(v, 10, 0); err == nil {
			mp[Int64Type] += 1
		} else if _, err := strconv.ParseFloat(v, 64); err == nil {
			mp[DoubleType] += 1
		} else {
			mp[StringType] += 1
		}
	}

	lenVs := len(vs)

	if mp[Int64Type] == lenVs {
		return Int64Type
	}

	if mp[Int64Type]+mp[DoubleType] == lenVs {
		return DoubleType
	}

	return StringType
}
