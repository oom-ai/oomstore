package oomstore

import (
	"context"
	"fmt"
	"io"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
)

func (s *OomStore) Apply(ctx context.Context, opt apply.ApplyOpt) error {
	stage, err := buildApplyStage(ctx, opt)
	if err != nil {
		return err
	}

	return s.metadata.WithTransaction(ctx, func(c context.Context, tx metadata.WriteStore) error {
		var (
			entityMp = make(map[string]int) // entity name -> entity id
			groupMp  = make(map[string]int) // group name -> group id
		)

		// apply entity
		for _, entity := range stage.NewEntities {
			if entityMp[entity.Name], err = s.applyEntity(ctx, tx, entity); err != nil {
				return err
			}
		}

		// apply group
		for _, group := range stage.NewGroups {
			if _, exist := entityMp[group.EntityName]; !exist {
				var entity *types.Entity
				if entity, err = s.metadata.GetEntityByName(c, group.EntityName); err != nil {
					return err
				}
				entityMp[group.EntityName] = entity.ID
			}

			group.EntityID = entityMp[group.EntityName]
			group.Category = types.BatchFeatureCategory
			if groupMp[group.Name], err = s.applyGroup(ctx, tx, group); err != nil {
				return err
			}
		}

		// apply feature
		for _, feature := range stage.NewFeatures {
			if _, exist := groupMp[feature.GroupName]; !exist {
				var group *types.Group
				if group, err = s.metadata.GetGroupByName(ctx, feature.GroupName); err != nil {
					return err
				}
				groupMp[feature.GroupName] = group.ID
			}

			feature.GroupID = groupMp[feature.GroupName]
			if err := s.applyFeature(ctx, tx, feature); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *OomStore) applyEntity(ctx context.Context, txStore metadata.WriteStore, newEntity apply.Entity) (int, error) {
	if err := newEntity.Validate(); err != nil {
		return 0, err
	}

	entity, err := s.metadata.GetEntityByName(ctx, newEntity.Name)
	if err != nil {
		if !errdefs.IsNotFound(err) {
			return 0, err
		}
		return txStore.CreateEntity(ctx, metadata.CreateEntityOpt{
			CreateEntityOpt: types.CreateEntityOpt{
				EntityName:  newEntity.Name,
				Length:      newEntity.Length,
				Description: newEntity.Description,
			},
		})
	}

	if newEntity.Description != entity.Description {
		err = txStore.UpdateEntity(ctx, metadata.UpdateEntityOpt{
			EntityID:       entity.ID,
			NewDescription: &newEntity.Description,
		})
		if err != nil {
			return 0, err
		}
	}

	return entity.ID, nil
}

func (s *OomStore) applyGroup(ctx context.Context, txStore metadata.WriteStore, newGroup apply.Group) (int, error) {
	if err := newGroup.Validate(); err != nil {
		return 0, err
	}

	group, err := s.metadata.GetGroupByName(ctx, newGroup.Name)
	if err != nil {
		if !errdefs.IsNotFound(err) {
			return 0, err
		}
		return txStore.CreateGroup(ctx, metadata.CreateGroupOpt{
			GroupName:   newGroup.Name,
			EntityID:    newGroup.EntityID,
			Category:    newGroup.Category,
			Description: newGroup.Description,
		})
	}

	if newGroup.Description != group.Description {
		err = txStore.UpdateGroup(ctx, metadata.UpdateGroupOpt{
			GroupID:        group.ID,
			NewDescription: &newGroup.Description,
		})
		if err != nil {
			return 0, err
		}

	}
	return group.ID, nil
}

func (s *OomStore) applyFeature(ctx context.Context, txStore metadata.WriteStore, newFeature apply.Feature) error {
	if err := newFeature.Validate(); err != nil {
		return err
	}

	feature, err := s.metadata.GetFeatureByName(ctx, newFeature.Name)
	if err != nil {
		if !errdefs.IsNotFound(err) {
			return err
		}
		valueType, err := s.offline.TypeTag(newFeature.DBValueType)
		if err != nil {
			return err
		}
		_, err = txStore.CreateFeature(ctx, metadata.CreateFeatureOpt{
			FeatureName: newFeature.Name,
			GroupID:     newFeature.GroupID,
			DBValueType: newFeature.DBValueType,
			Description: newFeature.Description,
			ValueType:   valueType,
		})
		return err
	}

	if newFeature.Description != feature.Description {
		return txStore.UpdateFeature(ctx, metadata.UpdateFeatureOpt{
			FeatureID:      feature.ID,
			NewDescription: &newFeature.Description,
		})
	}
	return nil
}

func buildApplyStage(ctx context.Context, opt apply.ApplyOpt) (*apply.ApplyStage, error) {
	var (
		stage = apply.NewApplyStage()
	)

	decoder := yaml.NewDecoder(opt.R)
	for {
		data := make(map[string]interface{})
		if err := decoder.Decode(data); err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		var kind string
		if k, ok := data["kind"]; ok {
			kind = k.(string)
		} else {
			return nil, fmt.Errorf("invalid yaml: missing kind")
		}

		switch kind {
		case "Entity":
			var entity apply.Entity
			if err := mapstructure.Decode(data, &entity); err != nil {
				return nil, err
			}

			// We don't want entity.BatchFeature to have values.
			// The whole stage should be a flat structure, not a nested structure.
			stage.NewEntities = append(stage.NewEntities, apply.Entity{
				Kind:        entity.Kind,
				Name:        entity.Name,
				Length:      entity.Length,
				Description: entity.Description,
			})

			for _, group := range entity.BatchFeatures {
				// We don't want group.Features to have values.
				// The whole stage should be a flat structure, not a nested structure.
				stage.NewGroups = append(stage.NewGroups, apply.Group{
					Kind:        "Group",
					Name:        group.Group,
					Group:       group.Group,
					EntityName:  entity.Name,
					Category:    group.Category,
					Description: group.Description,
				})

				for _, feature := range group.Features {
					stage.NewFeatures = append(stage.NewFeatures, apply.Feature{
						Kind:        "Feature",
						Name:        feature.Name,
						GroupName:   group.Group,
						DBValueType: feature.DBValueType,
						Description: feature.Description,
					})
				}
			}
		case "Group":
			var group apply.Group
			if err := mapstructure.Decode(data, &group); err != nil {
				return nil, err
			}

			// We don't want group.Features to have values.
			// The whole stage should be a flat structure, not a nested structure.
			stage.NewGroups = append(stage.NewGroups, apply.Group{
				Kind:        "Group",
				Name:        group.Name,
				Group:       group.Name,
				EntityName:  group.EntityName,
				Category:    group.Category,
				Description: group.Description,
			})

			for _, feature := range group.Features {
				stage.NewFeatures = append(stage.NewFeatures, apply.Feature{
					Kind:        "Feature",
					Name:        feature.Name,
					GroupName:   group.Name,
					DBValueType: feature.DBValueType,
					Description: feature.Description,
				})
			}

		case "Feature":
			var feature apply.Feature
			if err := mapstructure.Decode(data, &feature); err != nil {
				return nil, err
			}
			stage.NewFeatures = append(stage.NewFeatures, feature)

		default:
			return nil, fmt.Errorf("invalid kind '%s'", kind)
		}
	}
	return stage, nil
}
