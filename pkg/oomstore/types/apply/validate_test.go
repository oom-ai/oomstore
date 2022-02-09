package apply_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
)

func TestEntityValidate(t *testing.T) {
	cases := []struct {
		description string
		entity      *apply.Entity
		wantError   error
	}{
		{
			description: "invalid kind: spelling error",
			entity: &apply.Entity{
				Kind:        "Entit",
				Name:        "user",
				Description: "~",
			},
			wantError: fmt.Errorf(`Entity user: /kind: "Entit" must equal "Entity"`),
		},
		{
			description: "invalid kind: empty kind",
			entity: &apply.Entity{
				Name:        "user",
				Description: "~",
			},
			wantError: fmt.Errorf(`Entity user: /: {"Groups":null,"desc... "kind" value is required`),
		},
		{
			description: "invalid kind: all lowercase",
			entity: &apply.Entity{
				Kind:        "entity",
				Name:        "user",
				Description: "~",
			},
			wantError: fmt.Errorf(`Entity user: /kind: "entity" must equal "Entity"`),
		},
		{
			description: "empty name",
			entity: &apply.Entity{
				Kind:        "Entity",
				Description: "~",
			},
			wantError: fmt.Errorf(`Entity : /name: "" min length of 1 characters required: `),
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			err := c.entity.Validate()
			assert.Equal(t, c.wantError, err)
		})
	}
}

func TestGroupValidate(t *testing.T) {
	cases := []struct {
		description string
		group       *apply.Group
		wantError   error
	}{
		{
			description: "invalid kind: spelling error",
			group: &apply.Group{
				Kind:        "Grou",
				Name:        "student",
				Category:    types.CategoryBatch,
				EntityName:  "user",
				Description: "~",
			},
			wantError: fmt.Errorf(`Group student: /kind: "Grou" must equal "Group"`),
		},
		{
			description: "invalid kind: empty kind",
			group: &apply.Group{
				Name:        "student",
				Category:    types.CategoryBatch,
				EntityName:  "user",
				Description: "~",
			},
			wantError: fmt.Errorf(`Group student: /: {"category":"batch",... "kind" value is required`),
		},
		{
			description: "invalid kind: all lowercase",
			group: &apply.Group{
				Kind:        "group",
				Name:        "student",
				Category:    types.CategoryBatch,
				EntityName:  "user",
				Description: "~",
			},
			wantError: fmt.Errorf(`Group student: /kind: "group" must equal "Group"`),
		},
		{
			description: "empty name",
			group: &apply.Group{
				Kind:        "Group",
				Category:    types.CategoryBatch,
				EntityName:  "user",
				Description: "~",
			},
			wantError: fmt.Errorf(`Group : /name: "" min length of 1 characters required: `),
		},
		{
			description: "invalid category",
			group: &apply.Group{
				Kind:        "Group",
				Name:        "student",
				Category:    "batchh",
				EntityName:  "user",
				Description: "~",
			},
			wantError: fmt.Errorf(`Group student: /category: "batchh" should be one of ["batch", "stream"]`),
		},
		{
			description: "empty entity-name",
			group: &apply.Group{
				Kind:        "Group",
				Name:        "student",
				Category:    types.CategoryBatch,
				EntityName:  "",
				Description: "~",
			},
			wantError: fmt.Errorf(`Group student: /: {"category":"batch",... "entity-name" value is required`),
		},
		{
			description: "stream group snapshot-interval must greater than 0",
			group: &apply.Group{
				Kind:        "Group",
				Name:        "student",
				Category:    types.CategoryStream,
				EntityName:  "user",
				Description: "~",
			},
			wantError: fmt.Errorf(`Group student: /snapshot-interval: 0 0 must be greater than 0`),
		},
		{
			description: "batch group snapshot-interval must equal 0",
			group: &apply.Group{
				Kind:             "Group",
				Name:             "student",
				Category:         types.CategoryBatch,
				SnapshotInterval: time.Second,
				EntityName:       "user",
				Description:      "~",
			},
			wantError: fmt.Errorf(`Group student: /snapshot-interval: 1000000000 must equal 0`),
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			err := c.group.Validate()
			assert.Equal(t, c.wantError, err)
		})
	}
}

func TestFeatureValidate(t *testing.T) {
	cases := []struct {
		description string
		feature     *apply.Feature
		wantError   error
	}{
		{
			description: "invaid kind: spelling error",
			feature: &apply.Feature{
				Kind:        "Featur",
				Name:        "age",
				GroupName:   "student",
				ValueType:   "int64",
				Description: "~",
			},
			wantError: fmt.Errorf(`Feature age: /kind: "Featur" must equal "Feature"`),
		},
		{
			description: "invalid kind: empty kind ",
			feature: &apply.Feature{
				Name:        "age",
				GroupName:   "student",
				ValueType:   "int64",
				Description: "~",
			},
			wantError: fmt.Errorf(`Feature age: /: {"description":"~","... "kind" value is required`),
		},
		{
			description: "invalid kind: all lowercase",
			feature: &apply.Feature{
				Kind:        "feature",
				Name:        "age",
				GroupName:   "student",
				ValueType:   "int64",
				Description: "~",
			},
			wantError: fmt.Errorf(`Feature age: /kind: "feature" must equal "Feature"`),
		},
		{
			description: "empty group-name",
			feature: &apply.Feature{
				Kind:        "Feature",
				Name:        "age",
				ValueType:   "int64",
				Description: "~",
			},
			wantError: fmt.Errorf(`Feature age: /: {"description":"~","... "group-name" value is required`),
		},
		{

			description: "empty feature name",
			feature: &apply.Feature{
				Kind:        "Feature",
				GroupName:   "student",
				ValueType:   "int64",
				Description: "~",
			},
			wantError: fmt.Errorf(`Feature : /name: "" min length of 1 characters required: `),
		},
		{
			description: "invalid value type",
			feature: &apply.Feature{
				Kind:        "Feature",
				Name:        "age",
				GroupName:   "student",
				ValueType:   "int",
				Description: "~",
			},
			wantError: fmt.Errorf(`Feature age: /value-type: "int" should be one of ["string", "int64", "float64", "bool", "time", "bytes"]`),
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			err := c.feature.Validate()
			assert.Equal(t, c.wantError, err)
		})
	}
}
