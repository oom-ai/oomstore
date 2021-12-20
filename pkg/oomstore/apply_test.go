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
			description: "invalid yaml: missing kind or items",
			opt: apply.ApplyOpt{R: strings.NewReader(`
# kind: Entity
name: user
length: 8
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
		{
			description: "feature slice",
			opt: apply.ApplyOpt{
				R: strings.NewReader(`
items:
    - kind: Feature
      name: credit_score
      group-name: account
      db-value-type: int
      description: "credit_score description"
    - kind: Feature
      name: account_age_days
      group-name: account
      db-value-type: int
      description: "account_age_days description"
    - kind: Feature
      name: has_2fa_installed
      group-name: account
      db-value-type: bool
      description: "has_2fa_installed description"
    - kind: Feature
      name: transaction_count_7d
      group-name: transaction_stats
      db-value-type: int
      description: "transaction_count_7d description"
    - kind: Feature
      name: transaction_count_30d
      group-name: transaction_stats
      db-value-type: int
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
						DBValueType: "int",
						Description: "credit_score description",
					},
					{
						Kind:        "Feature",
						Name:        "account_age_days",
						GroupName:   "account",
						DBValueType: "int",
						Description: "account_age_days description",
					},
					{
						Kind:        "Feature",
						Name:        "has_2fa_installed",
						GroupName:   "account",
						DBValueType: "bool",
						Description: "has_2fa_installed description",
					},
					{
						Kind:        "Feature",
						Name:        "transaction_count_7d",
						GroupName:   "transaction_stats",
						DBValueType: "int",
						Description: "transaction_count_7d description",
					},
					{
						Kind:        "Feature",
						Name:        "transaction_count_30d",
						GroupName:   "transaction_stats",
						DBValueType: "int",
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
          db-value-type: int
          description: credit_score description
        - name: account_age_days
          db-value-type: int
          description: account_age_days description
        - name: has_2fa_installed
          db-value-type: bool
          description: has_2fa_installed description
    - kind: Group
      name: transaction_stats
      entity-name: user
      category: batch
      description: user transaction statistics
      features:
        - name: transaction_count_7d
          db-value-type: int
          description: transaction_count_7d description
        - name: transaction_count_30d
          db-value-type: int
          description: transaction_count_30d description
`),
			},
			wantStage: &apply.ApplyStage{
				NewEntities: []apply.Entity{},
				NewGroups: []apply.Group{
					{
						Kind:        "Group",
						Name:        "account",
						Group:       "account",
						EntityName:  "user",
						Category:    "batch",
						Description: "user account info",
					},
					{
						Kind:        "Group",
						Name:        "transaction_stats",
						Group:       "transaction_stats",
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
						DBValueType: "int",
						Description: "credit_score description",
					},
					{
						Kind:        "Feature",
						Name:        "account_age_days",
						GroupName:   "account",
						DBValueType: "int",
						Description: "account_age_days description",
					},
					{
						Kind:        "Feature",
						Name:        "has_2fa_installed",
						GroupName:   "account",
						DBValueType: "bool",
						Description: "has_2fa_installed description",
					},
					{
						Kind:        "Feature",
						Name:        "transaction_count_7d",
						GroupName:   "transaction_stats",
						DBValueType: "int",
						Description: "transaction_count_7d description",
					},
					{
						Kind:        "Feature",
						Name:        "transaction_count_30d",
						GroupName:   "transaction_stats",
						DBValueType: "int",
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
      length: 8
      description: user ID
      batch-features:
        - group: account
          description: user account info
          features:
            - name: credit_score
              db-value-type: int
              description: credit_score description
            - name: account_age_days
              db-value-type: int
              description: account_age_days description
            - name: has_2fa_installed
              db-value-type: bool
              description: has_2fa_installed description
        - group: transaction_stats
          description: user transaction statistics
          features:
            - name: transaction_count_7d
              db-value-type: int
              description: transaction_count_7d description
            - name: transaction_count_30d
              db-value-type: int
              description: transaction_count_30d description
    - kind: Entity
      name: device
      length: 8
      description: device info
      batch-features:
        - group: phone
          description: phone info
          features:
            - name: model
              db-value-type: varchar(32)
              description: model description
            - name: price
              db-value-type: int
              description: price description
`),
			},
			wantStage: &apply.ApplyStage{
				NewEntities: []apply.Entity{
					{
						Kind:        "Entity",
						Name:        "user",
						Length:      8,
						Description: "user ID",
					},
					{
						Kind:        "Entity",
						Name:        "device",
						Length:      8,
						Description: "device info",
					},
				},
				NewGroups: []apply.Group{
					{
						Kind:        "Group",
						Name:        "account",
						Group:       "account",
						EntityName:  "user",
						Category:    "batch",
						Description: "user account info",
					},
					{
						Kind:        "Group",
						Name:        "transaction_stats",
						Group:       "transaction_stats",
						EntityName:  "user",
						Category:    "batch",
						Description: "user transaction statistics",
					},
					{
						Kind:        "Group",
						Name:        "phone",
						Group:       "phone",
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
						DBValueType: "int",
						Description: "credit_score description",
					},
					{
						Kind:        "Feature",
						Name:        "account_age_days",
						GroupName:   "account",
						DBValueType: "int",
						Description: "account_age_days description",
					},
					{
						Kind:        "Feature",
						Name:        "has_2fa_installed",
						GroupName:   "account",
						DBValueType: "bool",
						Description: "has_2fa_installed description",
					},
					{
						Kind:        "Feature",
						Name:        "transaction_count_7d",
						GroupName:   "transaction_stats",
						DBValueType: "int",
						Description: "transaction_count_7d description",
					},
					{
						Kind:        "Feature",
						Name:        "transaction_count_30d",
						GroupName:   "transaction_stats",
						DBValueType: "int",
						Description: "transaction_count_30d description",
					},
					{
						Kind:        "Feature",
						Name:        "model",
						GroupName:   "phone",
						DBValueType: "varchar(32)",
						Description: "model description",
					},
					{
						Kind:        "Feature",
						Name:        "price",
						GroupName:   "phone",
						DBValueType: "int",
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
