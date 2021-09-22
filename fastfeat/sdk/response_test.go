package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FeatureValue(t *testing.T) {
	cases := []struct {
		msg string

		featureValue FeatureValue
		handle       func(FeatureValue) (interface{}, error)

		wantVal   interface{}
		wantError error
	}{
		{
			msg:          "case2: int feature type",
			featureValue: newFeatureValue("1"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.Int64()
			},

			wantVal:   int64(1),
			wantError: nil,
		},
		{
			msg:          "case3: float64 feature type",
			featureValue: newFeatureValue("3.1415926"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.Double()
			},

			wantVal:   3.1415926,
			wantError: nil,
		},
		{
			msg:          "case4: string feature type",
			featureValue: newFeatureValue("fastfeat response sdk"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.String()
			},

			wantVal:   "fastfeat response sdk",
			wantError: nil,
		},
		{
			msg:          "case5: int array feature type",
			featureValue: newFeatureValue("-1,0,1,2,3"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.Int64Array()
			},

			wantVal:   []int64{-1, 0, 1, 2, 3},
			wantError: nil,
		},
		{
			msg:          "case6: double array feature type",
			featureValue: newFeatureValue("1,2.0,3.1415926"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.DoubleArray()
			},

			wantVal:   []float64{1, 2, 3.1415926},
			wantError: nil,
		},
		{
			msg:          "case7: string array feature type",
			featureValue: newFeatureValue("伴鱼绘本,伴鱼英语,伴鱼数学"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.StringArray()
			},

			wantVal:   []string{"伴鱼绘本", "伴鱼英语", "伴鱼数学"},
			wantError: nil,
		},
		{
			msg:          "case9: err type mismatch",
			featureValue: newFeatureValue("伴鱼绘本"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.Int64()
			},

			wantVal:   int64(0),
			wantError: fmt.Errorf(ErrTypeMismatch, "int64"),
		},
		{
			msg:          "case10: err typ mismatch",
			featureValue: newFeatureValue("1,2,3,4"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.Int64()
			},

			wantVal:   int64(0),
			wantError: fmt.Errorf(ErrTypeMismatch, "int64"),
		},
		{
			// 思考：这种情况好像也说得通？暂时认为这是正确的，留作讨论
			msg:          "case12: int -> int array",
			featureValue: newFeatureValue("1"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.Int64Array()
			},

			wantVal:   []int64{1},
			wantError: nil,
		},
		{
			// 思考：这种情况好像也说得通？暂时认为这是正确的，留作讨论
			msg:          "case13: double -> double array",
			featureValue: newFeatureValue("1.0"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.DoubleArray()
			},

			wantVal:   []float64{1.0},
			wantError: nil,
		},
		{
			// 思考：这种情况好像也说得通？暂时认为这是正确的，留作讨论
			msg:          "case14: int -> double",
			featureValue: newFeatureValue("1"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.Double()
			},

			wantVal:   float64(1.0),
			wantError: nil,
		},
		{
			// 思考：这种情况好像也说得通？暂时认为这是正确的，留作讨论
			msg:          "case15: string -> string array",
			featureValue: newFeatureValue("伴鱼"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.StringArray()
			},

			wantVal:   []string{"伴鱼"},
			wantError: nil,
		},
		{
			// 思考：这种情况好像也说得通？暂时认为这是正确的，留作讨论
			msg:          "case16: int array -> double array",
			featureValue: newFeatureValue("1,2,3"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.DoubleArray()
			},

			wantVal:   []float64{1, 2, 3},
			wantError: nil,
		},
		{
			msg:          "case19: int array -> string array",
			featureValue: newFeatureValue("1,2,3,4,5"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.StringArray()
			},

			wantVal:   []string{"1", "2", "3", "4", "5"},
			wantError: nil,
		},
		{
			msg:          "case20: double array -> string array",
			featureValue: newFeatureValue("1.0,2,3.00"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.StringArray()
			},

			wantVal:   []string{"1.0", "2", "3.00"},
			wantError: nil,
		},
		{
			msg:          "case21: int(0) -> bool(false)",
			featureValue: newFeatureValue("0"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.Bool()
			},

			wantVal:   false,
			wantError: nil,
		},
		{
			msg:          "case22: int(1) -> bool(true)",
			featureValue: newFeatureValue("1"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.Bool()
			},

			wantVal:   true,
			wantError: nil,
		},
		{
			msg:          "case23: int(0/1/0) array -> bool(false/true/false) array",
			featureValue: newFeatureValue("0,1,0"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.BoolArray()
			},

			wantVal:   []bool{false, true, false},
			wantError: nil,
		},
		{
			msg:          "case24: int(1) -> double(1.0) array",
			featureValue: newFeatureValue("1"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.DoubleArray()
			},

			wantVal:   []float64{1.},
			wantError: nil,
		},
		{
			msg:          "case25: int(1) -> bool(true) array",
			featureValue: newFeatureValue("1"),
			handle: func(fv FeatureValue) (interface{}, error) {
				return fv.BoolArray()
			},

			wantVal:   []bool{true},
			wantError: nil,
		},
	}

	for _, c := range cases {
		val, err := c.handle(c.featureValue)
		assert.Equal(t, c.wantVal, val, c.msg)
		assert.Equal(t, c.wantError, err, c.msg)
	}
}
