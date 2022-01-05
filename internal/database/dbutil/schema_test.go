package dbutil_test

import (
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func TestBuildSchema(t *testing.T) {
	features := types.FeatureList{
		&types.Feature{
			Name:      "age",
			ValueType: types.Int64,
		},
		&types.Feature{
			Name:      "gender",
			ValueType: types.String,
		},
	}
	cases := []struct {
		description string

		tableName string
		entity    *types.Entity
		features  types.FeatureList
		backend   types.BackendType

		want    string
		wantErr error
	}{
		{
			description: "postgres schema",
			backend:     types.BackendPostgres,
			tableName:   "user",
			entity:      &types.Entity{Name: "user_id"},
			features:    features,

			want: `
CREATE TABLE "user" (
	"user_id" TEXT,
	"age" bigint,
	"gender" text
)`,
			wantErr: nil,
		},
		{
			description: "mysql schema",
			backend:     types.BackendMySQL,
			tableName:   "user",
			entity:      &types.Entity{Name: "user_id"},
			features:    features,

			want: "\n" +
				"CREATE TABLE `user` (\n" +
				"	`user_id` VARCHAR(255),\n" +
				"	`age` bigint,\n" +
				"	`gender` text\n)",
			wantErr: nil,
		},
		{
			description: "cassandra schema",
			backend:     types.BackendCassandra,
			tableName:   "user",
			entity: &types.Entity{
				Name: "user_id",
			},
			features: features,

			wantErr: nil,
			want: `
CREATE TABLE "user" (
	"user_id" TEXT,
	"age" bigint,
	"gender" text
)`,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			schema := dbutil.BuildTableSchema(c.tableName, c.entity, false, c.features, nil, c.backend)
			require.Equal(t, strings.TrimSpace(c.want), schema)
		})
	}
}
