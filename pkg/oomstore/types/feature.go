package types

import (
	"fmt"
	"time"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/pkg/oomstore/util"
)

type Feature struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	ValueType ValueType `db:"value_type"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`

	GroupID int `db:"group_id"`
	Group   *Group
}

func (f *Feature) FullName() string {
	if f.Group == nil {
		panic("expected group to be not nil")
	}
	return util.ComposeFullFeatureName(f.Group.Name, f.Name)
}

func (f *Feature) Copy() *Feature {
	if f == nil {
		return nil
	}
	copied := *f

	if copied.Group != nil {
		copied.Group = copied.Group.Copy()
	}
	return &copied
}

func (f *Feature) Entity() *Entity {
	return f.Group.Entity
}

func (f *Feature) DBFullName(backend BackendType) string {
	if backend != BackendBigQuery {
		return f.FullName()
	}
	return fmt.Sprintf("%s_%s", f.Group.Name, f.Name)
}

type FeatureList []*Feature

func (l FeatureList) Copy() FeatureList {
	if len(l) == 0 {
		return nil
	}
	copied := make(FeatureList, 0, len(l))
	for _, x := range l {
		copied = append(copied, x.Copy())
	}
	return copied
}

func (l *FeatureList) Len() int { return len(*l) }

func (l *FeatureList) Names() (names []string) {
	for _, f := range *l {
		names = append(names, f.Name)
	}
	return
}

func (l *FeatureList) FullNames() (fullNames []string) {
	for _, f := range *l {
		fullNames = append(fullNames, f.FullName())
	}
	return
}

func (l *FeatureList) IDs() (ids []int) {
	for _, f := range *l {
		ids = append(ids, f.ID)
	}
	return
}

func (l FeatureList) Filter(filter func(*Feature) bool) (rs FeatureList) {
	for _, f := range l {
		if filter(f) {
			rs = append(rs, f)
		}
	}
	return
}

func (l FeatureList) FilterFullNames(fullNames []string) (rs FeatureList) {
	return l.Filter(func(f *Feature) bool {
		for _, fullName := range fullNames {
			if f.FullName() == fullName {
				return true
			}
		}
		return false
	})
}

func (l FeatureList) Find(find func(*Feature) bool) *Feature {
	for _, f := range l {
		if find(f) {
			return f
		}
	}
	return nil
}

func (l FeatureList) GroupIDs() (ids []int) {
	groupIDMap := make(map[int]struct{})
	for _, r := range l {
		groupIDMap[r.GroupID] = struct{}{}
	}
	groupIDs := make([]int, 0, len(groupIDMap))
	for id := range groupIDMap {
		groupIDs = append(groupIDs, id)
	}
	return groupIDs
}

func (l FeatureList) GroupNames() []string {
	groupNameMap := make(map[string]struct{})
	groupNames := make([]string, 0, l.Len())
	for _, r := range l {
		if _, ok := groupNameMap[r.Group.Name]; !ok {
			groupNameMap[r.Group.Name] = struct{}{}
			groupNames = append(groupNames, r.Group.Name)
		}
	}
	return groupNames
}

func (l FeatureList) GroupByGroupID() map[int]FeatureList {
	featureMap := make(map[int]FeatureList)
	for _, f := range l {
		featureMap[f.GroupID] = append(featureMap[f.GroupID], f)
	}
	return featureMap
}

func (l FeatureList) GroupByGroupName() map[string]FeatureList {
	featureMap := make(map[string]FeatureList)
	for _, f := range l {
		featureMap[f.Group.Name] = append(featureMap[f.Group.Name], f)
	}
	return featureMap
}

func (l FeatureList) GetSharedEntity() (*Entity, error) {
	m := make(map[int]*Entity)
	for _, f := range l {
		m[f.Group.EntityID] = f.Group.Entity
	}
	if len(m) != 1 {
		return nil, errdefs.Errorf("expected 1 entity, got %d entities", len(m))
	}

	for _, entity := range m {
		return entity, nil
	}
	return nil, errdefs.Errorf("expected 1 entity, got 0")
}
