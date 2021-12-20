package dbutil_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func TestBuildSchema(t *testing.T) {
	features := types.FeatureList{
		&types.Feature{
			Name:      "age",
			ValueType: types.INT64,
		},
		&types.Feature{
			Name:      "gender",
			ValueType: types.STRING,
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
			backend:     types.POSTGRES,
			tableName:   "user",
			entity: &types.Entity{
				Name:   "user_id",
				Length: 32,
			},
			features: features,

			want: `
CREATE TABLE user (
	"user_id" VARCHAR(32) PRIMARY KEY,
	"age" bigint,
	"gender" text
)`,
			wantErr: nil,
		},
		{
			description: "mysql schema",
			backend:     types.MYSQL,
			tableName:   "user",
			entity: &types.Entity{
				Name:   "user_id",
				Length: 32,
			},
			features: features,

			want: "\n" +
				"CREATE TABLE user (\n" +
				"	`user_id` VARCHAR(32) PRIMARY KEY,\n" +
				"	`age` bigint,\n" +
				"	`gender` text\n)",
			wantErr: nil,
		},
		{
			description: "cassandra schema",
			backend:     types.CASSANDRA,
			tableName:   "user",
			entity: &types.Entity{
				Name: "user_id",
			},
			features: features,

			wantErr: nil,
			want: `
CREATE TABLE user (
	"user_id" TEXT PRIMARY KEY,
	"age" bigint,
	"gender" text
)`,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			schema, err := dbutil.BuildCreateSchema(c.tableName, c.entity, c.features, c.backend)
			require.Equal(t, c.wantErr, err)
			require.Equal(t, c.want, schema)
		})
	}
}
