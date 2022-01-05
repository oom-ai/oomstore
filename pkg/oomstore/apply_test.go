package oomstore

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
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
			description: "invalid yaml: missing kind or items",
			opt: apply.ApplyOpt{R: strings.NewReader(`
# kind: Entity
name: user
description: 'User ID'
`)},
			wantStage: nil,
			wantErr:   fmt.Errorf("invalid yaml: missing kind or items"),
		},
		{
			description: "invalid kind",
			opt: apply.ApplyOpt{R: strings.NewReader(`
kind: Entit
name: user
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
description: 'description'
`)},
			wantStage: &apply.ApplyStage{
				NewEntities: []apply.Entity{
					{
						Kind:        "Entity",
						Name:        "user",
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
value-type: string
description: 'description'
---
kind: Feature
name: price
group-name: device
category: batch
value-type: int64
description: 'description'
`)},
			wantStage: &apply.ApplyStage{
				NewEntities: []apply.Entity{
					{
						Kind:        "Entity",
						Name:        "user",
						Description: "description",
					},
				}, NewGroups: []apply.Group{
					{
						Kind:        "Group",
						Name:        "account",
						EntityName:  "user",
						Category:    "batch",
						Description: "description",
					},
					{
						Kind:        "Group",
						Name:        "device",
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
						ValueType:   "string",
						Description: "description",
					},
					{
						Kind:        "Feature",
						Name:        "price",
						GroupName:   "device",
						ValueType:   "int64",
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
  value-type: string
  description: 'description'
- name: price
  value-type: int64
  description: 'description'
`)},
			wantStage: &apply.ApplyStage{
				NewEntities: []apply.Entity{},
				NewGroups: []apply.Group{
					{
						Kind:        "Group",
						Name:        "device",
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
						ValueType:   "string",
						Description: "description",
					},
					{

						Kind:        "Feature",
						Name:        "price",
						GroupName:   "device",
						ValueType:   "int64",
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
description: 'description'
groups:
- name: device
  category: batch
  description: description
  features:
  - name: model
    value-type: string
    description: 'description'
  - name: price
    value-type: int64
    description: 'description'
- name: user
  category: batch
  description: description
  features:
  - name: age
    value-type: int64
    description: 'description'
  - name: gender
    value-type: int64
    description: 'description'
`)},
			wantStage: &apply.ApplyStage{
				NewEntities: []apply.Entity{
					{
						Kind:        "Entity",
						Name:        "user",
						Description: "description",
					},
				},
				NewGroups: []apply.Group{
					{
						Kind:        "Group",
						Name:        "device",
						Category:    types.CategoryBatch,
						EntityName:  "user",
						Description: "description",
					},
					{
						Kind:        "Group",
						Name:        "user",
						Category:    types.CategoryBatch,
						EntityName:  "user",
						Description: "description",
					},
				},
				NewFeatures: []apply.Feature{
					{
						Kind:        "Feature",
						Name:        "model",
						GroupName:   "device",
						ValueType:   "string",
						Description: "description",
					},
					{

						Kind:        "Feature",
						Name:        "price",
						GroupName:   "device",
						ValueType:   "int64",
						Description: "description",
					},
					{
						Kind:        "Feature",
						Name:        "age",
						GroupName:   "user",
						ValueType:   "int64",
						Description: "description",
					},
					{
						Kind:        "Feature",
						Name:        "gender",
						GroupName:   "user",
						ValueType:   "int64",
						Description: "description",
					},
				},
			},
		},
		{
			description: "feature slice",
			opt: apply.ApplyOpt{
				R: strings.NewReader(`
items:
    - kind: Feature
      name: credit_score
      group-name: account
      value-type: int64
      description: "credit_score description"
    - kind: Feature
      name: account_age_days
      group-name: account
      value-type: int64
      description: "account_age_days description"
    - kind: Feature
      name: has_2fa_installed
      group-name: account
      value-type: bool
      description: "has_2fa_installed description"
    - kind: Feature
      name: transaction_count_7d
      group-name: transaction_stats
      value-type: int64
      description: "transaction_count_7d description"
    - kind: Feature
      name: transaction_count_30d
      group-name: transaction_stats
      value-type: int64
      description: "transaction_count_30d description"
`),
			},
			wantStage: &apply.ApplyStage{
				NewEntities: []apply.Entity{},
				NewGroups:   []apply.Group{},
				NewFeatures: []apply.Feature{
					{
						Kind:        "Feature",
						Name:        "credit_score",
						GroupName:   "account",
						ValueType:   "int64",
						Description: "credit_score description",
					},
					{
						Kind:        "Feature",
						Name:        "account_age_days",
						GroupName:   "account",
						ValueType:   "int64",
						Description: "account_age_days description",
					},
					{
						Kind:        "Feature",
						Name:        "has_2fa_installed",
						GroupName:   "account",
						ValueType:   "bool",
						Description: "has_2fa_installed description",
					},
					{
						Kind:        "Feature",
						Name:        "transaction_count_7d",
						GroupName:   "transaction_stats",
						ValueType:   "int64",
						Description: "transaction_count_7d description",
					},
					{
						Kind:        "Feature",
						Name:        "transaction_count_30d",
						GroupName:   "transaction_stats",
						ValueType:   "int64",
						Description: "transaction_count_30d description",
					},
				},
			},
		},
		{
			description: "group slice",
			opt: apply.ApplyOpt{
				R: strings.NewReader(`
items:
    - kind: Group
      name: account
      entity-name: user
      category: batch
      description: user account info
      features:
        - name: credit_score
          value-type: int64
          description: credit_score description
        - name: account_age_days
          value-type: int64
          description: account_age_days description
        - name: has_2fa_installed
          value-type: bool
          description: has_2fa_installed description
    - kind: Group
      name: transaction_stats
      entity-name: user
      category: batch
      description: user transaction statistics
      features:
        - name: transaction_count_7d
          value-type: int64
          description: transaction_count_7d description
        - name: transaction_count_30d
          value-type: int64
          description: transaction_count_30d description
`),
			},
			wantStage: &apply.ApplyStage{
				NewEntities: []apply.Entity{},
				NewGroups: []apply.Group{
					{
						Kind:        "Group",
						Name:        "account",
						EntityName:  "user",
						Category:    "batch",
						Description: "user account info",
					},
					{
						Kind:        "Group",
						Name:        "transaction_stats",
						EntityName:  "user",
						Category:    "batch",
						Description: "user transaction statistics",
					},
				},
				NewFeatures: []apply.Feature{
					{
						Kind:        "Feature",
						Name:        "credit_score",
						GroupName:   "account",
						ValueType:   "int64",
						Description: "credit_score description",
					},
					{
						Kind:        "Feature",
						Name:        "account_age_days",
						GroupName:   "account",
						ValueType:   "int64",
						Description: "account_age_days description",
					},
					{
						Kind:        "Feature",
						Name:        "has_2fa_installed",
						GroupName:   "account",
						ValueType:   "bool",
						Description: "has_2fa_installed description",
					},
					{
						Kind:        "Feature",
						Name:        "transaction_count_7d",
						GroupName:   "transaction_stats",
						ValueType:   "int64",
						Description: "transaction_count_7d description",
					},
					{
						Kind:        "Feature",
						Name:        "transaction_count_30d",
						GroupName:   "transaction_stats",
						ValueType:   "int64",
						Description: "transaction_count_30d description",
					},
				},
			},
			wantErr: nil,
		},
		{
			description: "entity slice",
			opt: apply.ApplyOpt{
				R: strings.NewReader(`
items:
    - kind: Entity
      name: user
      description: user ID
      groups:
        - name: account
          category: batch
          description: user account info
          features:
            - name: credit_score
              value-type: int64
              description: credit_score description
            - name: account_age_days
              value-type: int64
              description: account_age_days description
            - name: has_2fa_installed
              value-type: bool
              description: has_2fa_installed description
        - name: transaction_stats
          category: batch
          description: user transaction statistics
          features:
            - name: transaction_count_7d
              value-type: int64
              description: transaction_count_7d description
            - name: transaction_count_30d
              value-type: int64
              description: transaction_count_30d description
    - kind: Entity
      name: device
      description: device info
      groups:
        - name: phone
          category: batch
          description: phone info
          features:
            - name: model
              value-type: string
              description: model description
            - name: price
              value-type: int64
              description: price description
`),
			},
			wantStage: &apply.ApplyStage{
				NewEntities: []apply.Entity{
					{
						Kind:        "Entity",
						Name:        "user",
						Description: "user ID",
					},
					{
						Kind:        "Entity",
						Name:        "device",
						Description: "device info",
					},
				},
				NewGroups: []apply.Group{
					{
						Kind:        "Group",
						Name:        "account",
						EntityName:  "user",
						Category:    "batch",
						Description: "user account info",
					},
					{
						Kind:        "Group",
						Name:        "transaction_stats",
						EntityName:  "user",
						Category:    "batch",
						Description: "user transaction statistics",
					},
					{
						Kind:        "Group",
						Name:        "phone",
						EntityName:  "device",
						Category:    "batch",
						Description: "phone info",
					},
				},
				NewFeatures: []apply.Feature{
					{
						Kind:        "Feature",
						Name:        "credit_score",
						GroupName:   "account",
						ValueType:   "int64",
						Description: "credit_score description",
					},
					{
						Kind:        "Feature",
						Name:        "account_age_days",
						GroupName:   "account",
						ValueType:   "int64",
						Description: "account_age_days description",
					},
					{
						Kind:        "Feature",
						Name:        "has_2fa_installed",
						GroupName:   "account",
						ValueType:   "bool",
						Description: "has_2fa_installed description",
					},
					{
						Kind:        "Feature",
						Name:        "transaction_count_7d",
						GroupName:   "transaction_stats",
						ValueType:   "int64",
						Description: "transaction_count_7d description",
					},
					{
						Kind:        "Feature",
						Name:        "transaction_count_30d",
						GroupName:   "transaction_stats",
						ValueType:   "int64",
						Description: "transaction_count_30d description",
					},
					{
						Kind:        "Feature",
						Name:        "model",
						GroupName:   "phone",
						ValueType:   "string",
						Description: "model description",
					},
					{
						Kind:        "Feature",
						Name:        "price",
						GroupName:   "phone",
						ValueType:   "int64",
						Description: "price description",
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
