package oomstore

import (
	"context"
	"io"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
)

// Apply registers entity, feature group, feature to metadata store.
func (s *OomStore) Apply(ctx context.Context, opt apply.ApplyOpt) error {
	stage, err := buildApplyStage(ctx, opt)
	if err != nil {
		return err
	}

	onlineJobs := make([]func() error, 0)
	if err := s.metadata.WithTransaction(ctx, func(c context.Context, tx metadata.DBStore) error {
		// apply entity
		for _, entity := range stage.NewEntities {
			if err2 := s.applyEntity(c, tx, entity); err2 != nil {
				return err2
			}
		}

		// apply group
		for _, group := range stage.NewGroups {
			onlineJob, err2 := s.applyGroup(c, tx, group)
			if err2 != nil {
				return err2
			}
			if onlineJob != nil {
				onlineJobs = append(onlineJobs, onlineJob)
			}
		}

		// apply feature
		for _, feature := range stage.NewFeatures {
			onlineJob, err2 := s.applyFeature(c, tx, feature)
			if err2 != nil {
				return err2
			}
			if onlineJob != nil {
				onlineJobs = append(onlineJobs, onlineJob)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	for _, job := range onlineJobs {
		if err := job(); err != nil {
			return err
		}
	}
	return nil
}

func (s *OomStore) applyEntity(ctx context.Context, tx metadata.DBStore, newEntity apply.Entity) error {
	if err := newEntity.Validate(); err != nil {
		return err
	}

	entity, err := tx.GetEntityByName(ctx, newEntity.Name)
	if err != nil {
		if !errdefs.IsNotFound(err) {
			return err
		}
		_, err = tx.CreateEntity(ctx, metadata.CreateEntityOpt{
			CreateEntityOpt: types.CreateEntityOpt{
				EntityName:  newEntity.Name,
				Description: newEntity.Description,
			},
		})
		return err
	}

	if newEntity.Description != entity.Description {
		return tx.UpdateEntity(ctx, metadata.UpdateEntityOpt{
			EntityID:       entity.ID,
			NewDescription: &newEntity.Description,
		})
	}
	return nil
}

func (s *OomStore) applyGroup(ctx context.Context, tx metadata.DBStore, newGroup apply.Group) (func() error, error) {
	if err := newGroup.Validate(); err != nil {
		return nil, err
	}

	entity, err := tx.GetEntityByName(ctx, newGroup.EntityName)
	if err != nil {
		return nil, err
	}

	group, err := tx.GetGroupByName(ctx, newGroup.Name)
	if err != nil {
		if !errdefs.IsNotFound(err) {
			return nil, err
		}

		_, err = tx.CreateGroup(ctx, metadata.CreateGroupOpt{
			GroupName:        newGroup.Name,
			EntityID:         entity.ID,
			Category:         newGroup.Category,
			SnapshotInterval: int(newGroup.SnapshotInterval.Seconds()),
			Description:      newGroup.Description,
		})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	if group.Entity.Name != newGroup.EntityName || group.Category != newGroup.Category {
		return nil, errdefs.Errorf("group %+v conflicts with already existing group %+v", newGroup, group)
	}

	opt := metadata.UpdateGroupOpt{
		GroupID: group.ID,
	}
	if newGroup.Description != group.Description {
		opt.NewDescription = &newGroup.Description
	}
	if newGroup.Category == types.CategoryStream && newGroup.SnapshotInterval != 0 {
		second := int(newGroup.SnapshotInterval.Seconds())
		opt.NewSnapshotInterval = &second
	}

	if opt.NewDescription != nil || opt.NewSnapshotInterval != nil {
		return nil, tx.UpdateGroup(ctx, metadata.UpdateGroupOpt{
			GroupID:             group.ID,
			NewSnapshotInterval: opt.NewSnapshotInterval,
			NewDescription:      &newGroup.Description,
		})
	}

	return nil, nil
}

func (s *OomStore) applyFeature(ctx context.Context, tx metadata.DBStore, newFeature apply.Feature) (func() error, error) {
	if err := newFeature.Validate(); err != nil {
		return nil, err
	}

	group, err := tx.GetGroupByName(ctx, newFeature.GroupName)
	if err != nil {
		return nil, err
	}

	feature, err := tx.GetFeatureByName(ctx, newFeature.GroupName, newFeature.Name)
	if err != nil {
		if !errdefs.IsNotFound(err) {
			return nil, err
		}

		revisions, err := tx.ListRevision(ctx, &group.ID)
		if err != nil {
			return nil, err
		}
		if len(revisions) > 0 {
			return nil, errdefs.Errorf("group %s already has data and cannot add features due to the join and export mechanism", group.Name)
		}

		valueType, err := types.ParseValueType(newFeature.ValueType)
		if err != nil {
			return nil, err
		}
		_, err = tx.CreateFeature(ctx, metadata.CreateFeatureOpt{
			FeatureName: newFeature.Name,
			GroupID:     group.ID,
			Description: newFeature.Description,
			ValueType:   valueType,
		})
		if err != nil {
			return nil, err
		}

		if group.Category == types.CategoryStream {
			features, err := tx.ListFeature(ctx, metadata.ListFeatureOpt{
				GroupIDs: &[]int{group.ID},
			})
			if err != nil {
				return nil, err
			}
			return func() error {
				return s.online.CreateTable(ctx, online.CreateTableOpt{
					EntityName: group.Entity.Name,
					TableName:  dbutil.OnlineStreamTableName(group.ID),
					Features:   features,
				})
			}, nil
		}
		return nil, nil
	}

	if newFeature.Description != feature.Description {
		return nil, tx.UpdateFeature(ctx, metadata.UpdateFeatureOpt{
			FeatureID:      feature.ID,
			NewDescription: &newFeature.Description,
		})
	}
	return nil, nil
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

		kind, err := parseKind(data)
		if err != nil {
			return nil, err
		}

		switch kind {
		case "Entity":
			var entity apply.Entity
			if err := mapstructureDecode(data, &entity); err != nil {
				return nil, errdefs.WithStack(err)
			}

			stage.NewEntities = append(stage.NewEntities, buildApplyEntity(entity))

			for _, group := range entity.Groups {
				stage.NewGroups = append(stage.NewGroups, buildApplyGroup(group, entity.Name))

				for _, feature := range group.Features {
					stage.NewFeatures = append(stage.NewFeatures, buildApplyFeature(feature, group.Name))
				}
			}
		case "Group":
			var group apply.Group
			if err := mapstructureDecode(data, &group); err != nil {
				return nil, errdefs.WithStack(err)
			}

			stage.NewGroups = append(stage.NewGroups, buildApplyGroup(group, group.EntityName))

			for _, feature := range group.Features {
				stage.NewFeatures = append(stage.NewFeatures, buildApplyFeature(feature, group.Name))
			}

		case "Feature":
			var feature apply.Feature
			if err := mapstructureDecode(data, &feature); err != nil {
				return nil, errdefs.WithStack(err)
			}
			stage.NewFeatures = append(stage.NewFeatures, feature)

		case "Items":
			itemsKind, err := parseItemsKind(data)
			if err != nil {
				return nil, err
			}
			switch itemsKind {
			case "Feature":
				featureItems := apply.FeatureItems{}
				if err := mapstructureDecode(data, &featureItems); err != nil {
					return nil, err
				}
				for _, item := range featureItems.Items {
					stage.NewFeatures = append(stage.NewFeatures, buildApplyFeature(item, item.GroupName))
				}
			case "Group":
				groupItems := apply.GroupItems{}
				if err := mapstructureDecode(data, &groupItems); err != nil {
					return nil, err
				}
				for _, group := range groupItems.Items {
					stage.NewGroups = append(stage.NewGroups, buildApplyGroup(group, group.EntityName))

					for _, feature := range group.Features {
						stage.NewFeatures = append(stage.NewFeatures, buildApplyFeature(feature, group.Name))
					}
				}
			case "Entity":
				entityItems := apply.EntityItems{}
				if err := mapstructureDecode(data, &entityItems); err != nil {
					return nil, err
				}

				for _, entity := range entityItems.Items {
					stage.NewEntities = append(stage.NewEntities, buildApplyEntity(entity))

					for _, group := range entity.Groups {
						stage.NewGroups = append(stage.NewGroups, buildApplyGroup(group, entity.Name))

						for _, feature := range group.Features {
							stage.NewFeatures = append(stage.NewFeatures, buildApplyFeature(feature, group.Name))
						}
					}
				}
			}
		default:
			return nil, errdefs.Errorf("invalid kind '%s'", kind)
		}
	}
	return stage, nil
}

func mapstructureDecode(data map[string]interface{}, result interface{}) error {
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     result,
	})
	if err != nil {
		return errdefs.WithStack(err)
	}
	if err := dec.Decode(data); err != nil {
		return errdefs.WithStack(err)
	}
	return nil
}

func parseKind(data map[string]interface{}) (string, error) {
	if k, ok := data["kind"]; ok {
		return cast.ToString(k), nil
	}
	if _, ok := data["items"]; ok {
		return "Items", nil
	}
	return "", errdefs.Errorf("invalid yaml: missing kind or items")
}

func parseItemsKind(data map[string]interface{}) (string, error) {
	items := apply.Items{}
	if err := mapstructureDecode(data, &items); err != nil {
		return "", err
	}
	return items.Kind(), nil
}

func buildApplyEntity(entity apply.Entity) apply.Entity {
	// We don't want entity.Groups to have values.
	// The whole stage should be a flat structure, not a nested structure.
	return apply.Entity{
		Kind:        entity.Kind,
		Name:        entity.Name,
		Description: entity.Description,
	}
}

func buildApplyGroup(group apply.Group, entityName string) apply.Group {
	// We don't want group.Features to have values.
	// The whole stage should be a flat structure, not a nested structure.
	return apply.Group{
		Kind:             "Group",
		Name:             group.Name,
		EntityName:       entityName,
		Category:         group.Category,
		SnapshotInterval: group.SnapshotInterval,
		Description:      group.Description,
	}
}

func buildApplyFeature(feature apply.Feature, groupName string) apply.Feature {
	return apply.Feature{
		Kind:        "Feature",
		Name:        feature.Name,
		GroupName:   groupName,
		ValueType:   feature.ValueType,
		Description: feature.Description,
	}
}
