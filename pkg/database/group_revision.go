package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type GroupRevision struct {
	Group       string    `db:"group"`
	Revision    string    `db:"revision"`
	Source      string    `db:"source"`
	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}

func (db *DB) ListGroupRevisionByGroup(ctx context.Context, group string) ([]GroupRevision, error) {
	query := "SELECT * FROM feature_revision AS fr WHERE fr.group = ?"
	revisions := make([]GroupRevision, 0)
	if err := db.SelectContext(ctx, &revisions, query, group); err != nil {
		return nil, err
	}
	return revisions, nil
}

func (r *GroupRevision) OneLineString() string {
	return strings.Join([]string{
		r.Group, r.Revision, r.Source, r.Description, r.CreateTime.Format(time.RFC3339), r.ModifyTime.Format(time.RFC3339)},
		",")
}

func RevisionExists(ctx context.Context, db *DB, group, revision string) error {
	_, err := getSourceTableNameByGroupAndRevision(ctx, db, group, revision)
	if err == sql.ErrNoRows {
		return fmt.Errorf("revision '%s' not found int feature group '%s'", revision, group)
	}
	return err
}

func getSourceTableNameByGroupAndRevision(ctx context.Context, db *DB, group string, revision string) (string, error) {
	var source string
	err := db.QueryRowContext(ctx,
		"select source from feature_revision where `group` = ? and revision = ?",
		group, revision).Scan(&source)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("revision not found: %s", revision)
	}
	return source, err
}

func GetFeatureValueType(ctx context.Context, db *DB, config *FeatureConfig) (string, error) {
	sourceTable, err := getSourceTableNameByGroupAndRevision(ctx, db, config.Group, config.Revision)
	if err != nil {
		return "", fmt.Errorf("failed fetching source table: %v", err)
	}

	column, err := db.ColumnInfo(ctx, sourceTable, config.Name)
	if err != nil {
		return "", fmt.Errorf("failed fetching source column: %v", err)
	}

	return column.Type, nil
}

func RegisterRevision(ctx context.Context, db *DB, group, revision, tableName, description string) error {
	_, err := db.ExecContext(ctx,
		"insert into feature_revision(`group`, revision, source, description) values(?, ?, ?, ?)",
		group,
		revision,
		tableName,
		description)
	return err
}
