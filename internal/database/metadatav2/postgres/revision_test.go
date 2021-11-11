package postgres_test

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/internal/database/metadatav2/postgres"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
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

	require.NoError(t, db.Refresh())
	groupName := "device_info"
	group, err := db.GetFeatureGroup(ctx, groupName)
	require.NoError(t, err)

	revision := typesv2.Revision{
		ID:        revisionId,
		Revision:  1000,
		GroupID:   groupId,
		DataTable: "device_info_1000",
		Anchored:  false,
		Group:     group,
	}

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
				ignoreCreateAndModifyTime(actual)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestListRevision(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	_, _, _, revisions := prepareRevisions(t, ctx, db)
	groupName := "device_info"
	var nilRevisionList typesv2.RevisionList
	require.NoError(t, db.Refresh())

	testCases := []struct {
		description string
		opt         metadatav2.ListRevisionOpt
		expected    typesv2.RevisionList
	}{
		{
			description: "list revision by groupName, succeed",
			opt: metadatav2.ListRevisionOpt{
				GroupName: &groupName,
			},
			expected: revisions,
		},
		{
			description: "list revision by dataTables, succeed",
			opt: metadatav2.ListRevisionOpt{
				DataTables: []string{"device_info_1000", "device_info_2000"},
			},
			expected: revisions,
		},
		{
			description: "list revision by invalid dataTables, return empty list",
			opt: metadatav2.ListRevisionOpt{
				DataTables: []string{"device_info_3000"},
			},
			expected: nilRevisionList,
		},
		{
			description: "list revision by empty dataTables, return empty list",
			opt: metadatav2.ListRevisionOpt{
				DataTables: []string{},
				GroupName:  &groupName,
			},
			expected: nilRevisionList,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := db.ListRevision(ctx, tc.opt)
			for _, item := range actual {
				ignoreCreateAndModifyTime(item)
			}
			sort.Slice(tc.expected, func(i, j int) bool {
				return tc.expected[i].ID < tc.expected[j].ID
			})
			sort.Slice(actual, func(i, j int) bool {
				return actual[i].ID < actual[j].ID
			})
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func boolPtr(i bool) *bool {
	return &i
}

func int32Ptr(i int32) *int32 {
	return &i
}

func ignoreCreateAndModifyTime(revision *typesv2.Revision) {
	revision.CreateTime = time.Time{}
	revision.ModifyTime = time.Time{}
}

func prepareRevisions(t *testing.T, ctx context.Context, db *postgres.DB) (int16, int16, []int32, typesv2.RevisionList) {
	entityID, err := db.CreateEntity(ctx, metadatav2.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	require.NoError(t, err)

	groupId, err := db.CreateFeatureGroup(ctx, metadatav2.CreateFeatureGroupOpt{
		Name:        "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	})
	require.NoError(t, err)

	revisionId1, err := db.CreateRevision(ctx, metadatav2.CreateRevisionOpt{
		Revision:  1000,
		GroupId:   groupId,
		DataTable: "device_info_1000",
		Anchored:  false,
	})
	require.NoError(t, err)

	revisionId2, err := db.CreateRevision(ctx, metadatav2.CreateRevisionOpt{
		Revision:  2000,
		GroupId:   groupId,
		DataTable: "device_info_2000",
		Anchored:  false,
	})
	require.NoError(t, err)

	require.NoError(t, db.Refresh())
	groupName := "device_info"
	group, err := db.GetFeatureGroup(ctx, groupName)
	require.NoError(t, err)

	revision1 := &typesv2.Revision{
		ID:        revisionId1,
		Revision:  1000,
		GroupID:   groupId,
		DataTable: "device_info_1000",
		Anchored:  false,
		Group:     group,
	}

	revision2 := &typesv2.Revision{
		ID:        revisionId2,
		Revision:  2000,
		GroupID:   groupId,
		DataTable: "device_info_2000",
		Anchored:  false,
		Group:     group,
	}

	return entityID, groupId, []int32{revisionId1, revisionId2}, typesv2.RevisionList{revision1, revision2}
}
