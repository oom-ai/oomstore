package oomstore_test

import (
	"context"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
	"github.com/stretchr/testify/require"
)

func TestApplyEntity(t *testing.T) {
	var entityYamlData = `
kind: Entity
name: user
length: 8
description: a description
batch-features:
- group: device
  description: a description
  features:
  - name: model
    db-type-value: varchar(16)
  - name: price
    db_type_value: int
-  group: user
   description: a description
   features:
   - name: age
     db-type-value: int
   - name: gender
     db-type-value: int
stream-features:
- name: c
  db-type-value: xxx
`
	store := oomstore.TEST__New(nil, nil, nil)

	require.Nil(t, store.Apply(context.Background(), apply.ApplyOpt{
		R: strings.NewReader(entityYamlData),
	}))
}
