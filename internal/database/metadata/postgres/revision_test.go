package postgres

import (
	"context"
	"fmt"
	"testing"

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
