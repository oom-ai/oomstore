package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type FeatureConfig struct {
	Name           string    `db:"name"`
	Group          string    `db:"group"`
	Revision       string    `db:"revision"`
	Status         string    `db:"status"`
	Category       string    `db:"category"`
	ValueType      string    `db:"value_type"`
	Description    string    `db:"description"`
	RevisionsLimit int       `db:"revisions_limit"`
	CreateTime     time.Time `db:"create_time"`
	ModifyTime     time.Time `db:"modify_time"`
}

func (db *DB) ListFeatureConfig(ctx context.Context) ([]FeatureConfig, error) {
	query := `SELECT * FROM feature_config`
	features := make([]FeatureConfig, 0)
	if err := db.SelectContext(ctx, &features, query); err != nil {
		return nil, err
	}
	return features, nil
}

func (db *DB) ListFeatureConfigByGroup(ctx context.Context, group string) ([]FeatureConfig, error) {
	query := "SELECT * FROM feature_config AS fc WHERE fc.group = ?"
	features := make([]FeatureConfig, 0)
	if err := db.SelectContext(ctx, &features, query, group); err != nil {
		return nil, err
	}
	return features, nil
}

func (db *DB) GetFeatureConfig(ctx context.Context, groupName, featureName string) (*FeatureConfig, error) {
	var feature FeatureConfig
	query := `SELECT * FROM feature_config AS fc WHERE fc.group = ? AND fc.name = ?`
	if err := db.GetContext(ctx, &feature, query, groupName, featureName); err != nil {
		return nil, err
	}
	return &feature, nil
}

func GetEntityTable(ctx context.Context, db *DB, group, featureName string) (string, error) {
	var revision string
	err := db.QueryRowContext(ctx, `select fc.revision from feature_config as fc where fc.group = ? and fc.name = ?`, group, featureName).Scan(&revision)
	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", err
	default:
		return group + "_" + revision, nil
	}
}

func UpdateFeatureConfig(ctx context.Context, db *DB, field string, value interface{}, group, name string) (sql.Result, error) {
	return db.ExecContext(ctx,
		fmt.Sprintf("update feature_config set %s = ? where `group` = ? and name = ?", field),
		value,
		group,
		name,
	)
}

func RegisterFeatureConfig(ctx context.Context, db *DB, config FeatureConfig) error {
	_, err := db.ExecContext(ctx,
		"insert into"+
			" feature_config(name, `group`, category, value_type, revision, revisions_limit, status, description)"+
			" values(?, ?, ?, ?, ?, ?, ?, ?)",
		config.Name,
		config.Group,
		config.Category,
		config.ValueType,
		config.Revision,
		config.RevisionsLimit,
		config.Status,
		config.Description,
	)
	return err
}

func (r *FeatureConfig) String() string {
	return strings.Join([]string{
		fmt.Sprintf("Name:           %s", r.Name),
		fmt.Sprintf("Group:          %s", r.Group),
		fmt.Sprintf("Revision:       %s", r.Revision),
		fmt.Sprintf("Status:         %s", r.Status),
		fmt.Sprintf("Category:       %s", r.Category),
		fmt.Sprintf("ValueType:      %s", r.ValueType),
		fmt.Sprintf("Description:    %s", r.Description),
		fmt.Sprintf("RevisionsLimit: %d", r.RevisionsLimit),
		fmt.Sprintf("CreateTime:     %s", r.CreateTime.Format(time.RFC3339)),
		fmt.Sprintf("ModifyTime:     %s", r.ModifyTime.Format(time.RFC3339)),
	}, "\n")
}

func (r *FeatureConfig) OneLineString() string {
	return strings.Join([]string{
		r.Name, r.Group, r.Revision, r.Status, r.Category, r.ValueType, r.Description,
		fmt.Sprintf("%d", r.RevisionsLimit),
		r.CreateTime.Format(time.RFC3339), r.ModifyTime.Format(time.RFC3339)},
		",")
}
