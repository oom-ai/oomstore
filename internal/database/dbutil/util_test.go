package dbutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
)

func TestPlaceholders(t *testing.T) {
	cases := []struct {
		describe string
		elem     string
		sep      string
		size     int
		want     string
	}{
		{
			describe: "one elem ?",

			elem: "?",
			sep:  ",",
			size: 1,
			want: "?",
		},
		{
			describe: "four elem ?",

			elem: "?",
			sep:  ",",
			size: 4,
			want: "?,?,?,?",
		},
		{
			describe: "one elem a",

			elem: "a",
			sep:  "-",
			size: 5,
			want: "a-a-a-a-a",
		},
		{
			describe: "empty elem",

			elem: "",
			sep:  ",",
			size: 4,
			want: ",,,",
		},
	}

	for _, c := range cases {
		t.Run(c.describe, func(t *testing.T) {
			got := dbutil.Placeholders(c.size, c.elem, c.sep)
			assert.Equal(t, c.want, got)
		})
	}
}
