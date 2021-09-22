package sdk

import (
	"fmt"
	"strings"
)

var (
	// ErrFeatureNotFound indicates that the a requested feature was not found in the response
	ErrFeatureNotFound = "Feature %s not found in response."

	// ErrTypeMismatch indicates that the there was a type mismatch in the returned values
	ErrTypeMismatch = "Requested output of type %s does not match type of feature value returned."
)

type FeatureValue struct {
	s string
}

var emptyFeatureValue = FeatureValue{s: FASTFEAT_EMPTY_VALUE}

func newFeatureValue(s string) FeatureValue {
	return FeatureValue{
		s: s,
	}
}

type GetOnlineFeaturesResponse struct {
        Features map[string]string
}

type FeatureWrapper map[string]FeatureValue

func NewFeatureWrapper(rawResponse GetOnlineFeaturesResponse) FeatureWrapper {
	if rawResponse.Features == nil {
		return FeatureWrapper{}
	}

	rs := make(FeatureWrapper, len(rawResponse.Features))

	for featureName, rawValue := range rawResponse.Features {
		rs[featureName] = newFeatureValue(rawValue)
	}
	return rs
}

func (f FeatureWrapper) FeatureValue(featureName string) FeatureValue {
	if f == nil {
		return emptyFeatureValue
	}
	if v, ok := f[featureName]; !ok {
		return emptyFeatureValue
	} else {
		return v
	}
}

func (v FeatureValue) String() (string, error) {
	switch value := parseRawFeatureResponse(v.s); value.(type) {
	case EmptyValue:
		return "", FastFeatEmptyValue
	default:
		return v.s, nil
	}
}

func (v FeatureValue) StringArray() ([]string, error) {
	switch value := parseRawFeatureResponse(v.s); value.(type) {
	case EmptyValue:
		return nil, FastFeatEmptyValue
	default:
		return strings.Split(v.s, ","), nil
	}
}

func (v FeatureValue) Int64() (int64, error) {
	switch value := parseRawFeatureResponse(v.s); value.(type) {
	case Int64Value:
		return value.(Int64Value).Int64(), nil
	case EmptyValue:
		return 0, FastFeatEmptyValue
	default:
		return 0, fmt.Errorf(ErrTypeMismatch, "int64")
	}
}

func (v FeatureValue) Int64Array() ([]int64, error) {
	switch value := parseRawFeatureResponse(v.s); value.(type) {
	case Int64Value:
		return []int64{value.(Int64Value).Int64()}, nil
	case Int64ArrayValue:
		return value.(Int64ArrayValue).Int64Array(), nil
	case EmptyValue:
		return nil, FastFeatEmptyValue
	default:
		return nil, fmt.Errorf(ErrTypeMismatch, "[]int64")
	}
}

func (v FeatureValue) Double() (float64, error) {
	switch value := parseRawFeatureResponse(v.s); value.(type) {
	case Int64Value:
		return float64(value.(Int64Value).Int64()), nil
	case DoubleValue:
		return value.(DoubleValue).Double(), nil
	case EmptyValue:
		return 0.0, FastFeatEmptyValue
	default:
		return 0.0, fmt.Errorf(ErrTypeMismatch, "double")
	}
}

func (v FeatureValue) DoubleArray() ([]float64, error) {
	switch value := parseRawFeatureResponse(v.s); value.(type) {
	case Int64Value:
		return []float64{float64(value.(Int64Value).Int64())}, nil
	case DoubleValue:
		return []float64{value.(DoubleValue).Double()}, nil
	case Int64ArrayValue:
		ivs := value.(Int64ArrayValue)
		rs := make([]float64, 0, len(ivs.Int64ArrayValue))

		for _, i := range ivs.Int64ArrayValue {
			rs = append(rs, float64(i.Int64Value))
		}
		return rs, nil
	case DoubleArrayValue:
		return value.(DoubleArrayValue).DoubleArray(), nil
	case EmptyValue:
		return nil, FastFeatEmptyValue
	default:
		return nil, fmt.Errorf(ErrTypeMismatch, "[]double")
	}
}

func (v FeatureValue) Bool() (bool, error) {
	switch value := parseRawFeatureResponse(v.s); value.(type) {
	// 将 1 和 0 解析为布尔值
	case Int64Value:
		v := value.(Int64Value).Int64()
		if v != 0 && v != 1 {
			return false, fmt.Errorf(ErrTypeMismatch, "bool")

		}
		return v == 1, nil
	case EmptyValue:
		return false, FastFeatEmptyValue
	default:
		return false, fmt.Errorf(ErrTypeMismatch, "bool")
	}
}

func (v FeatureValue) BoolArray() ([]bool, error) {
	switch value := parseRawFeatureResponse(v.s); value.(type) {
	case Int64Value:
		v := value.(Int64Value).Int64()
		if v != 0 && v != 1 {
			return nil, fmt.Errorf(ErrTypeMismatch, "[]bool")
		}
		return []bool{v == 1}, nil
	case Int64ArrayValue:
		var (
			vs = value.(Int64ArrayValue)
			rs = make([]bool, 0, len(vs.Int64ArrayValue))
		)
		for _, int64Value := range vs.Int64ArrayValue {
			v := int64Value.Int64()
			if v != 0 && v != 1 {
				return nil, fmt.Errorf(ErrTypeMismatch, "[]bool")
			}
			rs = append(rs, v == 1)
		}
		return rs, nil
	case EmptyValue:
		return nil, FastFeatEmptyValue
	default:
		return nil, fmt.Errorf(ErrTypeMismatch, "[]bool")
	}
}
