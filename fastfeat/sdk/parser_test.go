package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getRawFeatureResponseElementType(t *testing.T) {
	cases := []struct {
		msg         string
		rawResponse string
		want        FeatureElementTypeEnum
	}{
		{
			msg:         "case2: int type",
			rawResponse: "1,2,3",
			want:        Int64Type,
		},
		{
			msg:         "case3: double type",
			rawResponse: "1.00,3.1415926",
			want:        DoubleType,
		},
		{
			msg:         "case4: double type",
			rawResponse: "1,2,3.1415926",
			want:        DoubleType,
		},
		{
			msg:         "case5: string type",
			rawResponse: "hello world,nihao",
			want:        StringType,
		},
		{
			msg:         "case6: string type",
			rawResponse: "1,2,a",
			want:        StringType,
		},
		{
			msg:         "case7: string type",
			rawResponse: "1,2,3.1415926,true",
			want:        StringType,
		},
		{
			msg:         "case8: int type",
			rawResponse: "1",
			want:        Int64Type,
		},
	}

	for _, c := range cases {
		got := getRawFeatureResponseElementType(c.rawResponse)

		if got != c.want {
			t.Errorf("%s got %+v, but want %+v", c.msg, got, c.want)
		}
	}
}

func Test_parseFeatureValueResponse(t *testing.T) {
	cases := []struct {
		msg         string
		rawResponse string
		want        interface{}
	}{
		{
			msg:         "case 1: int type",
			rawResponse: "1",
			want:        Int64Value{1},
		},
		{
			msg:         "case 2: double type",
			rawResponse: "3.1415926",
			want:        DoubleValue{3.1415926},
		},
		{
			msg:         "case 5: string type",
			rawResponse: "this is plain text",
			want:        StringValue{"this is plain text"},
		},

		{
			msg:         "case 5: int slice type",
			rawResponse: "100,101,102,103,104",
			want: Int64ArrayValue{
				Int64ArrayValue: []Int64Value{
					{100},
					{101},
					{102},
					{103},
					{104},
				},
			},
		},
		{
			msg:         "case 7: double slice type",
			rawResponse: "1.00,3.1415926,2.0",
			want: DoubleArrayValue{
				DoubleArrayValue: []DoubleValue{
					{1.00},
					{3.1415926},
					{2.0},
				},
			},
		},
		{
			msg:         "case 8: double slice type two",
			rawResponse: "1,3.1415926,2.0",
			want: DoubleArrayValue{
				DoubleArrayValue: []DoubleValue{
					{1},
					{3.1415926},
					{2.0},
				},
			},
		},
		{
			msg:         "case 9: string slice type",
			rawResponse: "a,b,c,dd",
			want: StringArrayValue{
				StringArrayValue: []StringValue{
					{"a"},
					{"b"},
					{"c"},
					{"dd"},
				},
			},
		},
		{
			msg:         "case 10: string slice type",
			rawResponse: "1,2,3,a",
			want: StringArrayValue{
				StringArrayValue: []StringValue{
					{"1"},
					{"2"},
					{"3"},
					{"a"},
				},
			},
		},
		{
			msg:         "case 11: string slice type",
			rawResponse: "1,2,3,true",
			want: StringArrayValue{
				StringArrayValue: []StringValue{
					{"1"},
					{"2"},
					{"3"},
					{"true"},
				},
			},
		},
		{
			msg:         "case 12: string slice type",
			rawResponse: "1,3.1415926,false",
			want: StringArrayValue{
				StringArrayValue: []StringValue{
					{"1"},
					{"3.1415926"},
					{"false"},
				},
			},
		},
		{
			msg:         "case 13: string slice type",
			rawResponse: "1,2,3 ",
			want: StringArrayValue{
				StringArrayValue: []StringValue{
					{"1"},
					{"2"},
					{"3 "},
				},
			},
		},
	}

	for _, c := range cases {
		got := parseRawFeatureResponse(c.rawResponse)
		assert.Equal(t, c.want, got, c.msg)
	}
}

