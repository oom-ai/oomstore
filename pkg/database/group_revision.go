package database

import (
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

func ListGroupRevisionByGroup(db DB, group string) ([]GroupRevision, error) {
	query := "SELECT * FROM feature_revision AS fr WHERE fr.group = ?"
	revisions := make([]GroupRevision, 0)
	if err := db.Select(&revisions, query, group); err != nil {
		return nil, err
	}
	return revisions, nil
}

func (r *GroupRevision) OneLineString() string {
	return strings.Join([]string{
		r.Group, r.Revision, r.Source, r.Description, r.CreateTime.Format(time.RFC3339), r.ModifyTime.Format(time.RFC3339)},
		",")
}
