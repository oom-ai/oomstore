package oomstore

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
	"github.com/stretchr/testify/require"
)

func TestBuildApplyStage(t *testing.T) {
	testCases := []struct {
		description string
		opt         apply.ApplyOpt

		wantStage *apply.ApplyStage
		wantErr   error
	}{
		{
			description: "invalid yaml: missing kind",
			opt: apply.ApplyOpt{R: strings.NewReader(`
# kind: Entity
name: user
length: 8
description: 'User ID'
`)},
			wantStage: nil,
			wantErr:   fmt.Errorf("invalid yaml: missing kind"),
		},
		{
			description: "invalid kind",
			opt: apply.ApplyOpt{R: strings.NewReader(`
kind: Entit
name: user
length: 8
description: 'description'
`)},
			wantStage: nil,
			wantErr:   fmt.Errorf("invalid kind 'Entit'"),
		},
		{
			description: "single entity",
			opt: apply.ApplyOpt{R: strings.NewReader(`
kind: Entity
name: user
length: 8
description: 'description'
`)},
			wantStage: &apply.ApplyStage{
				NewEntities: []apply.Entity{
					{
						Kind:        "Entity",
						Name:        "user",
						Length:      8,
						Description: "description",
					},
				},
				NewGroups:   make([]apply.Group, 0),
				NewFeatures: make([]apply.Feature, 0),
			},
			wantErr: nil,
		},
		{
			description: "has many simple objects",
			opt: apply.ApplyOpt{R: strings.NewReader(`
kind: Entity
name: user
length: 8
description: 'description'
---
kind: Group
name: account
entity-name: user
category: batch
description: 'description'
---
kind: Group
name: device
entity-name: user
category: batch
description: 'description'
---
kind: Feature
name: model
group-name: device
category: batch
db-value-type: varchar(16)
description: 'description'
---
kind: Feature
name: price
group-name: device
category: batch
db-value-type: int
description: 'description'
`)},
			wantStage: &apply.ApplyStage{
				NewEntities: []apply.Entity{
					{
						Kind:        "Entity",
						Name:        "user",
						Length:      8,
						Description: "description",
					},
				}, NewGroups: []apply.Group{
					{
						Kind:        "Group",
						Name:        "account",
						Group:       "account",
						EntityName:  "user",
						Category:    "batch",
						Description: "description",
					},
					{
						Kind:        "Group",
						Name:        "device",
						Group:       "device",
						EntityName:  "user",
						Category:    "batch",
						Description: "description",
					},
				},
				NewFeatures: []apply.Feature{
					{
						Kind:        "Feature",
						Name:        "model",
						GroupName:   "device",
						DBValueType: "varchar(16)",
						Description: "description",
					},
					{
						Kind:        "Feature",
						Name:        "price",
						GroupName:   "device",
						DBValueType: "int",
						Description: "description",
					},
				},
			},
			wantErr: nil,
		},
		{
			description: "complex group",
			opt: apply.ApplyOpt{R: strings.NewReader(`
kind: Group
name: device
entity-name: user
category: batch
description: 'description'
features:
- name: model
  db-value-type: varchar(16)
  description: 'description'
- name: price
  db-value-type: int
  description: 'description'
`)},
			wantStage: &apply.ApplyStage{
				NewEntities: []apply.Entity{},
				NewGroups: []apply.Group{
					{
						Kind:        "Group",
						Name:        "device",
						Group:       "device",
						EntityName:  "user",
						Category:    "batch",
						Description: "description",
					},
				},
				NewFeatures: []apply.Feature{
					{
						Kind:        "Feature",
						Name:        "model",
						GroupName:   "device",
						DBValueType: "varchar(16)",
						Description: "description",
					},
					{

						Kind:        "Feature",
						Name:        "price",
						GroupName:   "device",
						DBValueType: "int",
						Description: "description",
					},
				},
			},
			wantErr: nil,
		},
		{
			description: "complex entity",
			opt: apply.ApplyOpt{R: strings.NewReader(`
kind: Entity
name: user
length: 8
description: 'description'
batch-features:
- group: device
  description: description
  features:
  - name: model
    db-value-type: varchar(16)
    description: 'description'
  - name: price
    db-value-type: int
    description: 'description'
- group: user
  description: description
  features:
  - name: age
    db-value-type: int
    description: 'description'
  - name: gender
    db-value-type: int
    description: 'description'
`)},
			wantStage: &apply.ApplyStage{
				NewEntities: []apply.Entity{
					{
						Kind:        "Entity",
						Name:        "user",
						Length:      8,
						Description: "description",
					},
				},
				NewGroups: []apply.Group{
					{
						Kind:        "Group",
						Name:        "device",
						Group:       "device",
						EntityName:  "user",
						Description: "description",
					},
					{
						Kind:        "Group",
						Name:        "user",
						Group:       "user",
						EntityName:  "user",
						Description: "description",
					},
				},
				NewFeatures: []apply.Feature{
					{
						Kind:        "Feature",
						Name:        "model",
						GroupName:   "device",
						DBValueType: "varchar(16)",
						Description: "description",
					},
					{

						Kind:        "Feature",
						Name:        "price",
						GroupName:   "device",
						DBValueType: "int",
						Description: "description",
					},
					{
						Kind:        "Feature",
						Name:        "age",
						GroupName:   "user",
						DBValueType: "int",
						Description: "description",
					},
					{
						Kind:        "Feature",
						Name:        "gender",
						GroupName:   "user",
						DBValueType: "int",
						Description: "description",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			stage, err := buildApplyStage(context.Background(), tc.opt)
			require.Equal(t, tc.wantErr, err)
			require.Equal(t, tc.wantStage, stage)
		})
	}
}
