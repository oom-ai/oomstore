package postgres_test

import (
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRevision(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	_, groupId := prepareEntityAndGroup(t, ctx, db)

	opt := metadatav2.CreateRevisionOpt{
		GroupId:     groupId,
		Revision:    1,
		DataTable:   "device_info_20211028",
		Description: "description",
	}

	testCases := []struct {
		description   string
		opt           metadatav2.CreateRevisionOpt
		expectedError error
		expected      int32
	}{
		{
			description:   "create revision successfully, return id",
			opt:           opt,
			expectedError: nil,
			expected:      int32(1),
		},
		{
			description:   "create existing revision, return error",
			opt:           opt,
			expectedError: fmt.Errorf("revision 1 already exist"),
			expected:      int32(0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := db.CreateRevision(ctx, tc.opt)
			assert.Equal(t, tc.expected, actual)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, tc.expectedError)
			}
		})
	}
}

func TestUpdateRevision(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	_, groupId := prepareEntityAndGroup(t, ctx, db)
	revisionId, err := db.CreateRevision(ctx, metadatav2.CreateRevisionOpt{
		Revision:  1000,
		GroupId:   groupId,
		DataTable: "device_info_1000",
		Anchored:  false,
	})
	require.NoError(t, err)

	testCases := []struct {
		description string
		opt         metadatav2.UpdateRevisionOpt
		expected    error
	}{
		{
			description: "update revision successfully",
			opt: metadatav2.UpdateRevisionOpt{
				RevisionID:  revisionId,
				NewAnchored: boolPtr(true),
			},
			expected: nil,
		},
		{
			description: "cannot update revision, return err",
			opt: metadatav2.UpdateRevisionOpt{
				RevisionID:  revisionId - 1,
				NewAnchored: boolPtr(true),
			},
			expected: fmt.Errorf("failed to update revision %d: revision not found", revisionId-1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := db.UpdateRevision(ctx, tc.opt)
			if tc.expected == nil {
				assert.NoError(t, actual)
			} else {
				assert.EqualError(t, actual, tc.expected.Error())
			}
		})
	}
}

func TestGetRevision(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	_, groupId := prepareEntityAndGroup(t, ctx, db)
	revisionId, err := db.CreateRevision(ctx, metadatav2.CreateRevisionOpt{
		Revision:  1000,
		GroupId:   groupId,
		DataTable: "device_info_1000",
		Anchored:  false,
	})
	require.NoError(t, err)
	revision := typesv2.Revision{
		ID:        revisionId,
		Revision:  1000,
		GroupID:   groupId,
		DataTable: "device_info_1000",
		Anchored:  false,
	}
	groupName := "device_info"
	db.Refresh()

	testCases := []struct {
		description   string
		opt           metadatav2.GetRevisionOpt
		expectedError error
		expected      *typesv2.Revision
	}{
		{
			description: "get revision by revisionId successfully",
			opt: metadatav2.GetRevisionOpt{
				RevisionId: &revisionId,
			},
			expectedError: nil,
			expected:      &revision,
		},
		{
			description: "get revision by groupName and revision successfully",
			opt: metadatav2.GetRevisionOpt{
				GroupName: &groupName,
				Revision:  &revision.Revision,
			},
			expectedError: nil,
			expected:      &revision,
		},
		{
			description: "get revision by groupName, return error",
			opt: metadatav2.GetRevisionOpt{
				GroupName: &groupName,
			},
			expectedError: fmt.Errorf("invalid GetRevisionOpt: %+v", metadatav2.GetRevisionOpt{
				GroupName: &groupName,
			}),
			expected: nil,
		},
		{
			description: "get revision by groupName, return error",
			opt: metadatav2.GetRevisionOpt{
				GroupName: &groupName,
			},
			expectedError: fmt.Errorf("invalid GetRevisionOpt: %+v", metadatav2.GetRevisionOpt{
				GroupName: &groupName,
			}),
			expected: nil,
		},
		{
			description: "get revision by revisionId and revision, return error",
			opt: metadatav2.GetRevisionOpt{
				RevisionId: &revisionId,
				Revision:   &revision.Revision,
			},
			expectedError: fmt.Errorf("invalid GetRevisionOpt: %+v", metadatav2.GetRevisionOpt{
				RevisionId: &revisionId,
				Revision:   &revision.Revision,
			}),
			expected: nil,
		},
		{
			description: "try to not existed revision, return error",
			opt: metadatav2.GetRevisionOpt{
				RevisionId: int32Ptr(0),
			},
			expectedError: fmt.Errorf("cannot find revision: revisionId=%d", 0),
			expected:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := db.GetRevision(ctx, tc.opt)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Equal(t, tc.expected, actual)
			} else {
				assert.NoError(t, tc.expectedError)
				tc.expected.CreateTime = actual.CreateTime
				tc.expected.ModifyTime = actual.ModifyTime
			}
		})
	}
}

func boolPtr(i bool) *bool {
	return &i
}

func int32Ptr(i int32) *int32 {
	return &i
}
