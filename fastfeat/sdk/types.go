package sdk

import "errors"

type FeatureElementTypeEnum string

const (
	StringType FeatureElementTypeEnum = "STRING"
	Int64Type  FeatureElementTypeEnum = "INT64"
	DoubleType FeatureElementTypeEnum = "DOUBLE"
)

const (
	FASTFEAT_EMPTY_VALUE = "FASTFEAT_EMPTY_VALUE"
)

var FastFeatEmptyValue = errors.New("fastfeat empty value")

type Value interface {
	IsValue()
}

type EmptyValue struct {
}

type ArrayValue struct {
	ArrayValue []Value
}

type StringValue struct {
	StringValue string
}

type Int64Value struct {
	Int64Value int64
}

type DoubleValue struct {
	DoubleValue float64
}

type StringArrayValue struct {
	StringArrayValue []StringValue
}

type Int64ArrayValue struct {
	Int64ArrayValue []Int64Value
}

type DoubleArrayValue struct {
	DoubleArrayValue []DoubleValue
}

func (EmptyValue) IsValue() {}

func (StringValue) IsValue() {}

func (s StringValue) String() string {
	return s.StringValue
}

func (StringArrayValue) IsValue() {}

func (sa StringArrayValue) StringArray() []string {
	rs := make([]string, 0, len(sa.StringArrayValue))
	for _, s := range sa.StringArrayValue {
		rs = append(rs, s.String())
	}
	return rs
}

func (Int64Value) IsValue() {}

func (i Int64Value) Int64() int64 {
	return i.Int64Value
}

func (Int64ArrayValue) IsValue() {}

func (ia Int64ArrayValue) Int64Array() []int64 {
	rs := make([]int64, 0, len(ia.Int64ArrayValue))
	for _, i := range ia.Int64ArrayValue {
		rs = append(rs, i.Int64())
	}
	return rs
}

func (DoubleValue) IsValue() {}

func (d DoubleValue) Double() float64 {
	return d.DoubleValue
}

func (DoubleArrayValue) IsValue() {}

func (da DoubleArrayValue) DoubleArray() []float64 {
	rs := make([]float64, 0, len(da.DoubleArrayValue))
	for _, d := range da.DoubleArrayValue {
		rs = append(rs, d.Double())
	}
	return rs
}

func (ArrayValue) IsValue() {}
