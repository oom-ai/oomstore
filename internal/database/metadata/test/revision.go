package test

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func TestCreateRevision(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	_, groupId := prepareEntityAndGroup(t, ctx, store)
	group, err := store.GetFeatureGroup(ctx, groupId)
	require.NoError(t, err)
	opt := metadata.CreateRevisionOpt{
		GroupID:     groupId,
		Revision:    1000,
		DataTable:   stringPtr("device_info_20211028"),
		Description: "description",
	}

	testCases := []struct {
		description      string
		opt              metadata.CreateRevisionOpt
		expectedError    error
		expected         int32
		expectedRevision *types.Revision
	}{
		{
			description:   "create revision successfully, return id",
			opt:           opt,
			expectedError: nil,
			expected:      int32(1),
			expectedRevision: &types.Revision{
				ID:          1,
				Revision:    1000,
				DataTable:   "device_info_20211028",
				Anchored:    false,
				Description: "description",
				GroupID:     groupId,
				Group:       group,
			},
		},
		{
			description: "create revision without data table, use default data table name",
			opt: metadata.CreateRevisionOpt{
				GroupID:     groupId,
				Revision:    2000,
				Description: "description",
			},
			expectedError: nil,
			expected:      int32(2),
			expectedRevision: &types.Revision{
				ID:          2,
				Revision:    2000,
				DataTable:   "data_1_2",
				Anchored:    false,
				Description: "description",
				GroupID:     groupId,
				Group:       group,
			},
		},
		{
			description:   "create existing revision, return error",
			opt:           opt,
			expectedError: fmt.Errorf("revision already exists: groupId=%d, revision=1000", groupId),
			expected:      int32(0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, _, err := store.CreateRevision(ctx, tc.opt)
			require.Equal(t, tc.expected, actual)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, tc.expectedError)
				require.NoError(t, store.Refresh())
				actualRevision, err := store.GetRevision(ctx, tc.expected)
				require.NoError(t, err)
				ignoreCreateAndModifyTime(actualRevision)
				require.Equal(t, tc.expectedRevision, actualRevision)
			}
		})
	}
}

func TestUpdateRevision(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	_, groupId := prepareEntityAndGroup(t, ctx, store)
	revisionId, _, err := store.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:  1000,
		GroupID:   groupId,
		DataTable: stringPtr("device_info_1000"),
		Anchored:  false,
	})
	require.NoError(t, err)

	testCases := []struct {
		description string
		opt         metadata.UpdateRevisionOpt
		expected    error
	}{
		{
			description: "update revision successfully",
			opt: metadata.UpdateRevisionOpt{
				RevisionID:  revisionId,
				NewAnchored: boolPtr(true),
			},
			expected: nil,
		},
		{
			description: "cannot update revision, return err",
			opt: metadata.UpdateRevisionOpt{
				RevisionID:  revisionId - 1,
				NewAnchored: boolPtr(true),
			},
			expected: fmt.Errorf("failed to update revision %d: revision not found", revisionId-1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := store.UpdateRevision(ctx, tc.opt)
			if tc.expected == nil {
				require.NoError(t, actual)
			} else {
				require.EqualError(t, actual, tc.expected.Error())
			}
		})
	}
}

func TestGetRevision(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	_, groupId := prepareEntityAndGroup(t, ctx, store)
	revisionId, _, err := store.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:  1000,
		GroupID:   groupId,
		DataTable: stringPtr("device_info_1000"),
		Anchored:  false,
	})
	require.NoError(t, err)

	require.NoError(t, store.Refresh())
	group, err := store.GetFeatureGroup(ctx, groupId)
	require.NoError(t, err)

	revision := types.Revision{
		ID:        revisionId,
		Revision:  1000,
		GroupID:   groupId,
		DataTable: "device_info_1000",
		Anchored:  false,
		Group:     group,
	}

	testCases := []struct {
		description   string
		revisionID    int32
		expectedError error
		expected      *types.Revision
	}{
		{
			description:   "get revision by revisionId successfully",
			revisionID:    revisionId,
			expectedError: nil,
			expected:      &revision,
		},
		{
			description:   "try to get not existed revision, return error",
			revisionID:    0,
			expectedError: fmt.Errorf("revision not found"),
			expected:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := store.GetRevision(ctx, tc.revisionID)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				require.Equal(t, tc.expected, actual)
			} else {
				require.NoError(t, tc.expectedError)
				ignoreCreateAndModifyTime(actual)
				require.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestGetRevisionBy(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	_, groupId := prepareEntityAndGroup(t, ctx, store)
	revisionId, _, err := store.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:  1000,
		GroupID:   groupId,
		DataTable: stringPtr("device_info_1000"),
		Anchored:  false,
	})
	require.NoError(t, err)

	require.NoError(t, store.Refresh())
	group, err := store.GetFeatureGroup(ctx, groupId)
	require.NoError(t, err)

	revision := types.Revision{
		ID:        revisionId,
		Revision:  1000,
		GroupID:   groupId,
		DataTable: "device_info_1000",
		Anchored:  false,
		Group:     group,
	}

	testCases := []struct {
		description   string
		opt           metadata.GetRevisionOpt
		GroupID       int16
		Revision      int64
		expectedError error
		expected      *types.Revision
	}{
		{
			description:   "get revision by groupID and revision successfully",
			GroupID:       groupId,
			Revision:      revision.Revision,
			expectedError: nil,
			expected:      &revision,
		},
		{
			description:   "try to get not existed revision, return error",
			GroupID:       groupId,
			Revision:      0,
			expectedError: fmt.Errorf("revision not found"),
			expected:      nil,
		},
		{
			description:   "try to get revision for a not existed group, return error",
			GroupID:       0,
			Revision:      revision.Revision,
			expectedError: fmt.Errorf("revision not found"),
			expected:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := store.GetRevisionBy(ctx, tc.GroupID, tc.Revision)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				require.Equal(t, tc.expected, actual)
			} else {
				require.NoError(t, tc.expectedError)
				ignoreCreateAndModifyTime(actual)
				require.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestListRevision(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	_, groupId, _, revisions := prepareRevisions(t, ctx, store)
	var nilRevisionList types.RevisionList
	require.NoError(t, store.Refresh())

	testCases := []struct {
		description string
		opt         metadata.ListRevisionOpt
		expected    types.RevisionList
	}{
		{
			description: "list revision by groupID, succeed",
			opt: metadata.ListRevisionOpt{
				GroupID: &groupId,
			},
			expected: revisions,
		},
		{
			description: "list revision by dataTables, succeed",
			opt: metadata.ListRevisionOpt{
				DataTables: []string{"device_info_1000", "device_info_2000"},
			},
			expected: revisions,
		},
		{
			description: "list revision by invalid dataTables, return empty list",
			opt: metadata.ListRevisionOpt{
				DataTables: []string{"device_info_3000"},
			},
			expected: nilRevisionList,
		},
		{
			description: "list revision by empty dataTables, return empty list",
			opt: metadata.ListRevisionOpt{
				DataTables: []string{},
				GroupID:    &groupId,
			},
			expected: nilRevisionList,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := store.ListRevision(ctx, tc.opt)
			for _, item := range actual {
				ignoreCreateAndModifyTime(item)
			}
			sort.Slice(tc.expected, func(i, j int) bool {
				return tc.expected[i].ID < tc.expected[j].ID
			})
			sort.Slice(actual, func(i, j int) bool {
				return actual[i].ID < actual[j].ID
			})
			require.Equal(t, tc.expected, actual)
		})
	}
}

func ignoreCreateAndModifyTime(revision *types.Revision) {
	revision.CreateTime = time.Time{}
	revision.ModifyTime = time.Time{}
}

func prepareRevisions(t *testing.T, ctx context.Context, store metadata.Store) (int16, int16, []int32, types.RevisionList) {
	entityID, err := store.CreateEntity(ctx, metadata.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	require.NoError(t, err)

	groupId, err := store.CreateFeatureGroup(ctx, metadata.CreateFeatureGroupOpt{
		Name:        "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	})
	require.NoError(t, err)
	require.NoError(t, store.Refresh())
	revisionId1, _, err := store.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:  1000,
		GroupID:   groupId,
		DataTable: stringPtr("device_info_1000"),
		Anchored:  false,
	})
	require.NoError(t, err)

	revisionId2, _, err := store.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:  2000,
		GroupID:   groupId,
		DataTable: stringPtr("device_info_2000"),
		Anchored:  false,
	})
	require.NoError(t, err)

	require.NoError(t, store.Refresh())
	group, err := store.GetFeatureGroup(ctx, groupId)
	require.NoError(t, err)

	revision1 := &types.Revision{
		ID:        revisionId1,
		Revision:  1000,
		GroupID:   groupId,
		DataTable: "device_info_1000",
		Anchored:  false,
		Group:     group,
	}

	revision2 := &types.Revision{
		ID:        revisionId2,
		Revision:  2000,
		GroupID:   groupId,
		DataTable: "device_info_2000",
		Anchored:  false,
		Group:     group,
	}

	return entityID, groupId, []int32{revisionId1, revisionId2}, types.RevisionList{revision1, revision2}
}
