package postgres

import (
	"testing"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func TestValueTypeTag(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected string
	}{
		{
			"bit varying",
			types.BYTE_ARRAY,
		},
		{
			"int",
			types.INT32,
		},
		{
			"varchar(32)",
			types.STRING,
		},
		{
			"timestamp without time zone",
			types.TIME,
		},
	} {
		actual, err := TypeTag(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		if actual != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, actual)
		}
	}
}
