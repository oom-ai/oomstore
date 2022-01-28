package cmd

import (
	"time"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type FlattenEntity struct {
	ID          int    `oomcli:"ID"`
	Name        string `oomcli:"NAME"`
	Description string `oomcli:"DESCRIPTION,truncate"`

	CreateTime time.Time `oomcli:"CREATE-TIME,wide"`
	ModifyTime time.Time `oomcli:"MODIFY-TIME,wide"`
}

type FlattenGroup struct {
	ID               int            `oomcli:"ID"`
	Name             string         `oomcli:"NAME"`
	Entity           string         `oomcli:"ENTITY"`
	Category         types.Category `oomcli:"CATEGORY"`
	SnapshotInterval time.Duration  `oomcli:"SNAPSHOT-INTERVAL"`
	Description      string         `oomcli:"DESCRIPTION,truncate"`

	OnlineRevisionID *int      `oomcli:"ONLINE-REVISION-ID,wide"`
	CreateTime       time.Time `oomcli:"CREATE-TIME,wide"`
	ModifyTime       time.Time `oomcli:"MODIFY-TIME,wide"`
}

type FlattenFeature struct {
	ID          int             `oomcli:"ID"`
	Name        string          `oomcli:"NAME"`
	Group       string          `oomcli:"GROUP"`
	Entity      string          `oomcli:"ENTITY"`
	Category    types.Category  `oomcli:"CATEGORY"`
	ValueType   types.ValueType `oomcli:"VALUE-TYPE"`
	Description string          `oomcli:"DESCRIPTION,truncate"`

	OnlineRevisionID *int      `oomcli:"ONLINE-REVISION-ID,wide"`
	CreateTime       time.Time `oomcli:"CREATE-TIME,wide"`
	ModifyTime       time.Time `oomcli:"MODIFY-TIME,wide"`
}

type FlattenRevision struct {
	ID            int    `oomcli:"ID"`
	Revision      int64  `oomcli:"REVISION"`
	Group         string `oomcli:"GROUP"`
	SnapshotTable string `oomcli:"SNAPSHOT-TABLE"`
	CdcTable      string `oomcli:"CDC-TABLE"`
	Description   string `oomcli:"DESCRIPTION,truncate"`

	Anchored   bool      `oomcli:"ANCHORED,wide"`
	CreateTime time.Time `oomcli:"CREATE-TIME,wide"`
	ModifyTime time.Time `oomcli:"MODIFY-TIME,wide"`
}

func parseTokenLists(i interface{}) (headerTokens TokenList, dataTokens []TokenList, err error) {
	var element interface{}
	var tokens TokenList
	switch s := i.(type) {
	case types.EntityList:
		element = FlattenEntity{}
		for _, e := range s {
			tokens, err = parseTokens(FlattenEntity{
				ID:          e.ID,
				Name:        e.Name,
				Description: e.Description,
				CreateTime:  e.CreateTime,
				ModifyTime:  e.ModifyTime,
			})
			if err != nil {
				return
			}
			dataTokens = append(dataTokens, tokens)
		}
	case types.FeatureList:
		element = FlattenFeature{}
		for _, e := range s {
			tokens, err = parseTokens(FlattenFeature{
				ID:               e.ID,
				Name:             e.Name,
				Group:            e.Group.Name,
				Entity:           e.Entity().Name,
				Category:         e.Group.Category,
				ValueType:        e.ValueType,
				Description:      e.Description,
				OnlineRevisionID: e.Group.OnlineRevisionID,
				CreateTime:       e.CreateTime,
				ModifyTime:       e.ModifyTime,
			})
			if err != nil {
				return
			}
			dataTokens = append(dataTokens, tokens)
		}
	case types.GroupList:
		element = FlattenGroup{}
		for _, e := range s {
			tokens, err = parseTokens(FlattenGroup{
				ID:               e.ID,
				Name:             e.Name,
				Entity:           e.Entity.Name,
				Category:         e.Category,
				SnapshotInterval: time.Duration(e.SnapshotInterval) * time.Second,
				Description:      e.Description,
				OnlineRevisionID: e.OnlineRevisionID,
				CreateTime:       e.CreateTime,
				ModifyTime:       e.ModifyTime,
			})
			if err != nil {
				return
			}
			dataTokens = append(dataTokens, tokens)
		}
	case types.RevisionList:
		element = FlattenRevision{}
		for _, e := range s {
			tokens, err = parseTokens(FlattenRevision{
				ID:            e.ID,
				Revision:      e.Revision,
				Group:         e.Group.Name,
				SnapshotTable: e.SnapshotTable,
				Description:   e.Description,
				Anchored:      e.Anchored,
				CreateTime:    e.CreateTime,
				ModifyTime:    e.ModifyTime,
			})
			if err != nil {
				return
			}
			dataTokens = append(dataTokens, tokens)
		}
	default:
		return nil, nil, errdefs.Errorf("unsupported type %T", i)
	}
	headerTokens, err = parseTokens(element)
	if err != nil {
		return
	}
	return
}
