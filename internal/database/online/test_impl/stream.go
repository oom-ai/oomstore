package test_impl

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type StreamData struct {
	entity   types.Entity
	groups   types.GroupList
	features types.FeatureList
}

var simpleStreamData StreamData

func init() {
	entity := types.Entity{
		ID:   1,
		Name: "user",
	}

	group1 := &types.Group{
		ID:       1,
		Name:     "user_clicks",
		Category: types.CategoryStream,
		EntityID: 1,
		Entity:   &entity,
	}

	group2 := &types.Group{
		ID:       2,
		Name:     "user_reads",
		Category: types.CategoryStream,
		EntityID: 1,
		Entity:   &entity,
	}

	simpleStreamData = StreamData{
		entity: entity,
		groups: types.GroupList{group1, group2},
		features: types.FeatureList{
			&types.Feature{
				ID:        1,
				Name:      "last_5_click_post",
				ValueType: types.String,
				GroupID:   1,
				Group:     group1,
			},
			&types.Feature{
				ID:        2,
				Name:      "last_10_click_post",
				ValueType: types.String,
				GroupID:   1,
				Group:     group1,
			},
			&types.Feature{
				ID:        3,
				Name:      "recent_5_read_post",
				ValueType: types.String,
				GroupID:   2,
				Group:     group2,
			},
		},
	}
}
