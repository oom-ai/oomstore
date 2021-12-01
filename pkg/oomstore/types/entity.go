package types

import (
	"fmt"
	"strings"
	"time"
)

type Entity struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	Length int    `db:"length"`

	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	ModifyTime  time.Time `db:"modify_time"`
}

func (e *Entity) Copy() *Entity {
	if e == nil {
		return nil
	}
	copied := *e
	return &copied
}

type EntityList []*Entity

func (l *EntityList) Find(find func(*Entity) bool) *Entity {
	for _, e := range *l {
		if find(e) {
			return e
		}
	}
	return nil
}

func (e *Entity) String() string {
	return strings.Join([]string{
		fmt.Sprintf("Name:        %s", e.Name),
		fmt.Sprintf("Length:      %d", e.Length),
		fmt.Sprintf("Description: %s", e.Description),
		fmt.Sprintf("CreateTime:  %s", e.CreateTime.Format(time.RFC3339)),
		fmt.Sprintf("ModifyTime:  %s", e.ModifyTime.Format(time.RFC3339)),
	}, "\n")
}
