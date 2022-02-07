package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestRevisionListBefore(t *testing.T) {
	testCases := []struct {
		description  string
		revisionList types.RevisionList
		unixMilli    int64
		expected     *types.Revision
	}{
		{
			description:  "empty revision list, return nil",
			revisionList: nil,
			unixMilli:    10,
			expected:     nil,
		},
		{
			description: "one revision",
			revisionList: []*types.Revision{
				{Revision: 10},
			},
			unixMilli: 11,
			expected:  &types.Revision{Revision: 10},
		},
		{
			description: "unixMilli less than the smallest revision, return nil",
			revisionList: []*types.Revision{
				{Revision: 10},
				{Revision: 5},
			},
			unixMilli: 2,
			expected:  nil,
		},
		{
			description: "unixMilli equal to the smallest revision, return revision",
			revisionList: []*types.Revision{
				{Revision: 10},
				{Revision: 5},
			},
			unixMilli: 5,
			expected:  &types.Revision{Revision: 5},
		},
		{
			description: "unixMilli greater than the largest revision, return the largest revision",
			revisionList: []*types.Revision{
				{Revision: 10},
				{Revision: 5},
				{Revision: 7},
			},
			unixMilli: 15,
			expected:  &types.Revision{Revision: 10},
		},
		{
			description: "unixMilli equal to the largest revision, return the largest revision",
			revisionList: []*types.Revision{
				{Revision: 10},
				{Revision: 5},
				{Revision: 7},
			},
			unixMilli: 10,
			expected:  &types.Revision{Revision: 10},
		},
		{
			description: "unixMilli greater than to the middle revision, return the middle revision",
			revisionList: []*types.Revision{
				{Revision: 10},
				{Revision: 5},
				{Revision: 7},
			},
			unixMilli: 8,
			expected:  &types.Revision{Revision: 7},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := tc.revisionList.Before(tc.unixMilli)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
