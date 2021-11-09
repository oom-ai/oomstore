package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/stretchr/testify/assert"
)

func TestCreateRevision(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	opt := metadata.CreateRevisionOpt{
		GroupName:   "device_baseinfo",
		Revision:    1,
		DataTable:   "device_bastinfo_20211028",
		Description: "description",
	}

	revision, err := db.CreateRevision(context.Background(), opt)
	assert.Nil(t, err)
	assert.Equal(t, int32(1), revision.ID)

	revision, err = db.CreateRevision(context.Background(), opt)
	assert.Equal(t, err, fmt.Errorf("revision 1 already exist"))
	assert.Nil(t, revision)
}

func TestUpdateRevision(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	groupName := "device_baseinfo"
	revisionTimestamp := time.Now().Unix()
	opt := metadata.CreateRevisionOpt{
		GroupName:   groupName,
		Revision:    revisionTimestamp,
		DataTable:   "device_bastinfo_20211028",
		Description: "description",
	}
	revision, err := db.CreateRevision(context.Background(), opt)
	assert.Nil(t, err)
	assert.Equal(t, int32(1), revision.ID)

	r, err := db.GetRevision(context.Background(), metadata.GetRevisionOpt{
		GroupName: &groupName,
		Revision:  &revisionTimestamp,
	})
	assert.NoError(t, err)

	newRevison := time.Now().Add(time.Second).Unix()
	rowsAffected, err := db.UpdateRevision(context.Background(), metadata.UpdateRevisionOpt{
		RevisionID:  int64(r.ID),
		NewRevision: &newRevison,
	})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	r, err = db.GetRevision(context.Background(), metadata.GetRevisionOpt{
		RevisionId: &r.ID,
	})
	assert.Nil(t, err)
	assert.Equal(t, newRevison, r.Revision)
}

func TestUpdateRevisionWithEmpty(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	newRevision := int64(0)
	rowsAffected, err := db.UpdateRevision(context.Background(), metadata.UpdateRevisionOpt{
		RevisionID:  0,
		NewRevision: &newRevision,
	})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), rowsAffected)
}

func TestGetRevision(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	r, err := db.GetRevision(context.Background(), metadata.GetRevisionOpt{})
	assert.NotNil(t, err)
	assert.Nil(t, r)

	opt := metadata.CreateRevisionOpt{
		GroupName:   "device_baseinfo",
		Revision:    1,
		DataTable:   "device_bastinfo_20211028",
		Description: "description",
	}
	revision, err := db.CreateRevision(context.Background(), opt)
	assert.Nil(t, err)
	assert.Equal(t, int32(1), revision.ID)

	groupName := "invalid-group-name"
	r, err = db.GetRevision(context.Background(), metadata.GetRevisionOpt{
		GroupName: &groupName,
	})
	assert.NotNil(t, err)
	assert.Nil(t, r)

	r, err = db.GetRevision(context.Background(), metadata.GetRevisionOpt{})
	assert.Nil(t, err)
	assert.Equal(t, opt.GroupName, r.GroupName)
	assert.Equal(t, opt.Revision, r.Revision)
	assert.Equal(t, opt.DataTable, r.DataTable)
	assert.Equal(t, opt.Description, r.Description)
}

func TestListRevision(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	rs, err := db.ListRevision(context.Background(), metadata.ListRevisionOpt{})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(rs))

	opt1 := metadata.CreateRevisionOpt{
		GroupName:   "device_baseinfo",
		Revision:    1,
		DataTable:   "device_bastinfo_20211028",
		Description: "description",
	}

	opt2 := metadata.CreateRevisionOpt{
		GroupName:   "device_baseinfo",
		Revision:    2,
		DataTable:   "device_bastinfo_20211029",
		Description: "description",
	}

	revision, err := db.CreateRevision(context.Background(), opt1)
	assert.Nil(t, err)
	assert.Equal(t, int32(1), revision.ID)

	revision, err = db.CreateRevision(context.Background(), opt2)
	assert.Nil(t, err)
	assert.Equal(t, int32(2), revision.ID)

	groupName := "device_baseinfo"
	rs, err = db.ListRevision(context.Background(), metadata.ListRevisionOpt{
		GroupName: &groupName,
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rs))

	rs, err = db.ListRevision(context.Background(), metadata.ListRevisionOpt{
		DataTables: []string{},
	})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(rs))

	rs, err = db.ListRevision(context.Background(), metadata.ListRevisionOpt{
		DataTables: []string{"device_bastinfo_20211028"},
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(rs))
}
