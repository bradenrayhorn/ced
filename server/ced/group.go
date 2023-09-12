package ced

import (
	"context"
	"strings"
)

type Group struct {
	ID           ID
	Name         Name
	Attendees    uint8
	MaxAttendees uint8
	HasResponded bool
	SearchHints  string
}

func NewGroup(name Name, maxAttendees uint8, searchHints string) (Group, error) {
	if err := ValidateFields(Field("Name", name)); err != nil {
		return Group{}, err
	}

	group := Group{
		ID:           NewID(),
		Name:         Name(strings.TrimSpace(string(name))),
		MaxAttendees: maxAttendees,
		Attendees:    0,
		HasResponded: false,
		SearchHints:  strings.TrimSpace(searchHints),
	}

	return group, nil
}

type GroupImport struct {
	Name         Name
	MaxAttendees uint8
	SearchHints  string
}

type GroupContract interface {
	// Creates a new group.
	Create(ctx context.Context, name Name, maxAttendees uint8, searchHints string) (Group, error)

	// Searches for a group using the search string.
	Search(ctx context.Context, search string) ([]Group, error)

	// Gets a group by id.
	Get(ctx context.Context, id ID) (Group, error)

	// Updates response for a group.
	Respond(ctx context.Context, id ID, attendees uint8, connectingIP string) error

	// Imports a list of groups.
	Import(ctx context.Context, groups []GroupImport) error
}

type GroupRespository interface {
	Create(ctx context.Context, group Group) error
	CreateMany(ctx context.Context, groups []Group) error
	Update(ctx context.Context, group Group) error
	Get(ctx context.Context, id ID) (Group, error)

	SearchByName(ctx context.Context, search string) ([]Group, error)
}
