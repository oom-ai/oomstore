package typesv2

import (
	"fmt"
	"strings"
	"time"
)

type FeatureGroup struct {
	ID       int16  `db:"id"`
	Name     string `db:"name"`
	Category string `db:"category"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`

	EntityID         int16  `db:"entity_id"`
	OnlineRevisionID *int32 `db:"online_revision_id"`

	Entity         *Entity
	OnlineRevision *Revision
}

func (fg *FeatureGroup) Copy() *FeatureGroup {
	return fg.copyWith(nil)
}

func (fg *FeatureGroup) copyWith(onlineRevision *Revision) *FeatureGroup {
	if fg == nil {
		return nil
	}
	copied := *fg

	if copied.OnlineRevisionID != nil {
		id := *copied.OnlineRevisionID
		copied.OnlineRevisionID = &id
	}

	if onlineRevision != nil {
		copied.OnlineRevision = onlineRevision
	} else if copied.OnlineRevision != nil {
		revision := copied.OnlineRevision.copyWith(&copied)
		copied.OnlineRevision = revision
	}

	if copied.Entity != nil {
		entity := copied.Entity.Copy()
		copied.Entity = entity
	}

	return &copied
}

type FeatureGroupList []*FeatureGroup

func (l *FeatureGroupList) Find(find func(*FeatureGroup) bool) *FeatureGroup {
	for _, g := range *l {
		if find(g) {
			return g
		}
	}
	return nil
}

func (l *FeatureGroupList) Filter(filter func(*FeatureGroup) bool) (rs FeatureGroupList) {
	for _, g := range *l {
		if filter(g) {
			rs = append(rs, g)
		}
	}
	return
}

func (fg *FeatureGroup) String() string {
	onlineRevision := "<NULL>"

	if fg.OnlineRevision != nil {
		onlineRevision = fmt.Sprint(*fg.OnlineRevision)
	}
	return strings.Join([]string{
		fmt.Sprintf("Name:            %s", fg.Name),
		fmt.Sprintf("Entity:          %s", fg.Entity.Name),
		fmt.Sprintf("Description:     %s", fg.Description),
		fmt.Sprintf("Online Revision: %s", onlineRevision),
		fmt.Sprintf("CreateTime:      %s", fg.CreateTime.Format(time.RFC3339)),
		fmt.Sprintf("ModifyTime:      %s", fg.ModifyTime.Format(time.RFC3339)),
	}, "\n")
}
