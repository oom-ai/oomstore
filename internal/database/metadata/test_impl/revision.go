package test_impl

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

	_, groupID := prepareEntityAndGroup(t, ctx, store)
	group, err := store.GetGroup(ctx, groupID)
	require.NoError(t, err)
	opt := metadata.CreateRevisionOpt{
		GroupID:     groupID,
		Revision:    1000,
		DataTable:   stringPtr("device_info_20211028"),
		Description: "description",
	}

	testCases := []struct {
		description      string
		opt              metadata.CreateRevisionOpt
		expectedError    error
		expected         int
		expectedRevision *types.Revision
	}{
		{
			description:   "create revision successfully, return id",
			opt:           opt,
			expectedError: nil,
			expected:      1,
			expectedRevision: &types.Revision{
				ID:          1,
				Revision:    1000,
				DataTable:   "device_info_20211028",
				Anchored:    false,
				Description: "description",
				GroupID:     groupID,
				Group:       group,
			},
		},
		{
			description: "create revision without data table, use default data table name",
			opt: metadata.CreateRevisionOpt{
				GroupID:     groupID,
				Revision:    2000,
				Description: "description",
			},
			expectedError: nil,
			expected:      2,
			expectedRevision: &types.Revision{
				ID:          2,
				Revision:    2000,
				DataTable:   "offline_1_2",
				Anchored:    false,
				Description: "description",
				GroupID:     groupID,
				Group:       group,
			},
		},
		{
			description:   "create existing revision, return error",
			opt:           opt,
			expectedError: fmt.Errorf("revision already exists: groupID=%d, revision=1000", groupID),
			expected:      0,
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
				actualRevision, err := store.CacheGetRevision(ctx, tc.expected)
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

	_, groupID := prepareEntityAndGroup(t, ctx, store)
	revisionID, _, err := store.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:  1000,
		GroupID:   groupID,
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
				RevisionID:  revisionID,
				NewAnchored: boolPtr(true),
			},
			expected: nil,
		},
		{
			description: "cannot update revision, return err",
			opt: metadata.UpdateRevisionOpt{
				RevisionID:  revisionID - 1,
				NewAnchored: boolPtr(true),
			},
			expected: fmt.Errorf("failed to update revision %d: revision not found", revisionID-1),
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

	_, groupID := prepareEntityAndGroup(t, ctx, store)
	revisionID, _, err := store.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:  1000,
		GroupID:   groupID,
		DataTable: stringPtr("device_info_1000"),
		Anchored:  false,
	})
	require.NoError(t, err)

	require.NoError(t, store.Refresh())
	group, err := store.GetGroup(ctx, groupID)
	require.NoError(t, err)

	revision := types.Revision{
		ID:        revisionID,
		Revision:  1000,
		GroupID:   groupID,
		DataTable: "device_info_1000",
		Anchored:  false,
		Group:     group,
	}

	testCases := []struct {
		description   string
		revisionID    int
		expectedError error
		expected      *types.Revision
	}{
		{
			description:   "get revision by revisionID successfully",
			revisionID:    revisionID,
			expectedError: nil,
			expected:      &revision,
		},
		{
			description:   "try to get not existed revision, return error",
			revisionID:    0,
			expectedError: fmt.Errorf("revision 0 not found"),
			expected:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := store.CacheGetRevision(ctx, tc.revisionID)

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

	_, groupID := prepareEntityAndGroup(t, ctx, store)
	revisionID, _, err := store.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:  1000,
		GroupID:   groupID,
		DataTable: stringPtr("device_info_1000"),
		Anchored:  false,
	})
	require.NoError(t, err)

	require.NoError(t, store.Refresh())
	group, err := store.GetGroup(ctx, groupID)
	require.NoError(t, err)

	revision := types.Revision{
		ID:        revisionID,
		Revision:  1000,
		GroupID:   groupID,
		DataTable: "device_info_1000",
		Anchored:  false,
		Group:     group,
	}

	testCases := []struct {
		description   string
		GroupID       int
		Revision      int64
		expectedError error
		expected      *types.Revision
	}{
		{
			description:   "get revision by groupID and revision successfully",
			GroupID:       groupID,
			Revision:      revision.Revision,
			expectedError: nil,
			expected:      &revision,
		},
		{
			description:   "try to get not existed revision, return error",
			GroupID:       groupID,
			Revision:      0,
			expectedError: fmt.Errorf("revision not found by group %d and revision 0", groupID),
			expected:      nil,
		},
		{
			description:   "try to get revision for a not existed group, return error",
			GroupID:       0,
			Revision:      revision.Revision,
			expectedError: fmt.Errorf("revision not found by group 0 and revision %d", revision.Revision),
			expected:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := store.CacheGetRevisionBy(ctx, tc.GroupID, tc.Revision)

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

	_, groupID, _, revisions := prepareRevisions(t, ctx, store)
	require.NoError(t, store.Refresh())

	testCases := []struct {
		description string
		groupID     *int
		expected    types.RevisionList
	}{
		{
			description: "list revision, succeed",
			groupID:     nil,
			expected:    revisions,
		},
		{
			description: "list revision by groupID, succeed",
			groupID:     &groupID,
			expected:    revisions,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := store.CacheListRevision(ctx, tc.groupID)
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

func prepareRevisions(t *testing.T, ctx context.Context, store metadata.Store) (int, int, []int, types.RevisionList) {
	entityID, err := store.CreateEntity(ctx, metadata.CreateEntityOpt{
		CreateEntityOpt: types.CreateEntityOpt{
			EntityName:  "device",
			Length:      32,
			Description: "description",
		},
	})
	require.NoError(t, err)

	groupID, err := store.CreateGroup(ctx, metadata.CreateGroupOpt{
		GroupName:   "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	})
	require.NoError(t, err)
	require.NoError(t, store.Refresh())
	revisionID1, _, err := store.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:  1000,
		GroupID:   groupID,
		DataTable: stringPtr("device_info_1000"),
		Anchored:  false,
	})
	require.NoError(t, err)

	revisionID2, _, err := store.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:  2000,
		GroupID:   groupID,
		DataTable: stringPtr("device_info_2000"),
		Anchored:  false,
	})
	require.NoError(t, err)

	require.NoError(t, store.Refresh())
	group, err := store.GetGroup(ctx, groupID)
	require.NoError(t, err)

	revision1 := &types.Revision{
		ID:        revisionID1,
		Revision:  1000,
		GroupID:   groupID,
		DataTable: "device_info_1000",
		Anchored:  false,
		Group:     group,
	}

	revision2 := &types.Revision{
		ID:        revisionID2,
		Revision:  2000,
		GroupID:   groupID,
		DataTable: "device_info_2000",
		Anchored:  false,
		Group:     group,
	}

	return entityID, groupID, []int{revisionID1, revisionID2}, types.RevisionList{revision1, revision2}
}
