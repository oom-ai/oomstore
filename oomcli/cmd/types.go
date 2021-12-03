package cmd

import (
	"fmt"
	"time"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type FlattenEntity struct {
	ID          int    `oomcli:"ID"`
	Name        string `oomcli:"NAME"`
	Length      int    `oomcli:"LENGTH"`
	Description string `oomcli:"DESCRIPTION,truncate"`

	CreateTime time.Time `oomcli:"CREATE-TIME,wide"`
	ModifyTime time.Time `oomcli:"MODIFY-TIME,wide"`
}

type FlattenGroup struct {
	ID          int    `oomcli:"ID"`
	Name        string `oomcli:"NAME"`
	Entity      string `oomcli:"ENTITY"`
	Description string `oomcli:"DESCRIPTION,truncate"`

	OnlineRevisionID *int      `oomcli:"ONLINE-REVISION-ID,wide"`
	CreateTime       time.Time `oomcli:"CREATE-TIME,wide"`
	ModifyTime       time.Time `oomcli:"MODIFY-TIME,wide"`
}

type FlattenFeature struct {
	ID          int    `oomcli:"ID"`
	Name        string `oomcli:"NAME"`
	Group       string `oomcli:"GROUP"`
	Entity      string `oomcli:"ENTITY"`
	Category    string `oomcli:"CATEGORY"`
	ValueType   string `oomcli:"VALUE-TYPE"`
	Description string `oomcli:"DESCRIPTION,truncate"`

	DBValueType      string    `oomcli:"DB-VALUE-TYPE,wide"`
	OnlineRevisionID *int      `oomcli:"ONLINE-REVISION-ID,wide"`
	CreateTime       time.Time `oomcli:"CREATE-TIME,wide"`
	ModifyTime       time.Time `oomcli:"MODIFY-TIME,wide"`
}

type FlattenRevision struct {
	ID          int    `oomcli:"ID"`
	Revision    int64  `oomcli:"REVISION"`
	Group       string `oomcli:"GROUP"`
	DataTable   string `oomcli:"DATA-TABLE"`
	Description string `oomcli:"DESCRIPTION,truncate"`

	Anchored   bool      `oomcli:"ANCHORED,wide"`
	CreateTime time.Time `oomcli:"CREATE-TIME,wide"`
	ModifyTime time.Time `oomcli:"MODIFY-TIME,wide"`
}

func parseTokenLists(i interface{}) (rs []TokenList, err error) {
	switch s := i.(type) {
	case types.EntityList:
		for _, e := range s {
			tokens, err := parseTokens(FlattenEntity{
				ID:          e.ID,
				Name:        e.Name,
				Length:      e.Length,
				Description: e.Description,
				CreateTime:  e.CreateTime,
				ModifyTime:  e.ModifyTime,
			})
			if err != nil {
				return nil, err
			}
			rs = append(rs, tokens)
		}
		return
	case types.FeatureList:
		for _, e := range s {
			tokens, err := parseTokens(FlattenFeature{
				ID:               e.ID,
				Name:             e.Name,
				Group:            e.Group.Name,
				Entity:           e.Entity().Name,
				Category:         e.Group.Category,
				ValueType:        e.ValueType,
				Description:      e.Description,
				DBValueType:      e.DBValueType,
				OnlineRevisionID: e.OnlineRevisionID(),
				CreateTime:       e.CreateTime,
				ModifyTime:       e.ModifyTime,
			})
			if err != nil {
				return nil, err
			}
			rs = append(rs, tokens)
		}
		return
	case types.GroupList:
		for _, e := range s {
			tokens, err := parseTokens(FlattenGroup{
				ID:               e.ID,
				Name:             e.Name,
				Entity:           e.Entity.Name,
				Description:      e.Description,
				OnlineRevisionID: e.OnlineRevisionID,
				CreateTime:       e.CreateTime,
				ModifyTime:       e.ModifyTime,
			})
			if err != nil {
				return nil, err
			}
			rs = append(rs, tokens)
		}
		return
	case types.RevisionList:
		for _, e := range s {
			tokens, err := parseTokens(FlattenRevision{
				ID:          e.ID,
				Revision:    e.Revision,
				Group:       e.Group.Name,
				DataTable:   e.DataTable,
				Description: e.Description,
				Anchored:    e.Anchored,
				CreateTime:  e.CreateTime,
				ModifyTime:  e.ModifyTime,
			})
			if err != nil {
				return nil, err
			}
			rs = append(rs, tokens)
		}
		return
	default:
		return nil, fmt.Errorf("unsupported type %T", i)
	}
}
