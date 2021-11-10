package postgres_test

import (
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRevision(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	entityId, err := db.CreateEntity(ctx, types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "device entity",
	})
	require.NoError(t, err)
	groupId, err := db.CreateFeatureGroup(ctx, metadatav2.CreateFeatureGroupOpt{
		Name:        "device_info",
		EntityID:    entityId,
		Category:    types.BatchFeatureCategory,
		Description: "device info",
	})
	require.NoError(t, err)

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
			if tc.expectedError != nil {
				assert.Error(t, err, tc.expectedError)
				assert.Equal(t, tc.expected, actual)
			} else {
				assert.NoError(t, tc.expectedError)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
