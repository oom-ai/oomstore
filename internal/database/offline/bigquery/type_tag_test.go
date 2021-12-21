package bigquery_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline/bigquery"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestValueTypeTag(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected string
	}{
		{
			"Numeric(10,1)",
			types.FLOAT64,
		},
		{
			"SMALLINT",
			types.INT64,
		},
		{
			"STRING(10)",
			types.STRING,
		},
		{
			"datetime",
			types.TIME,
		},
	} {
		actual, err := bigquery.TypeTag(tc.input)
		if err != nil {
			t.Fatal(err)
		}

		if actual != tc.expected {
			t.Errorf("expected %s, got %s", tc.expected, actual)
		}
	}
}
