package ced

import (
	"context"
	"fmt"
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

func (g Group) Validate() error {
	if err := ValidateFields(Field("Name", g.Name)); err != nil {
		return err
	}

	if g.Attendees > g.MaxAttendees {
		return NewError(EINVALID, fmt.Sprintf("group can have at most %d attendees, has %d", g.MaxAttendees, g.Attendees))
	}
	return nil
}

func NewGroup(name Name, maxAttendees uint8, searchHints string) (Group, error) {
	group := Group{
		ID:           NewID(),
		Name:         Name(strings.TrimSpace(string(name))),
		MaxAttendees: maxAttendees,
		Attendees:    0,
		HasResponded: false,
		SearchHints:  strings.TrimSpace(searchHints),
	}

	if err := group.Validate(); err != nil {
		return Group{}, err
	}

	return group, nil
}

type GroupImport struct {
	Name         Name
	MaxAttendees uint8
	SearchHints  string
}

type GroupUpdate struct {
	ID           ID
	Name         *Name
	Attendees    *uint8
	MaxAttendees *uint8
	SearchHints  *string
}

type ReqContext struct {
	ConnectingIP string
}

type GroupContract interface {
	// Creates a new group.
	Create(ctx context.Context, name Name, maxAttendees uint8, searchHints string) (Group, error)

	// Searches for a group using the search string.
	Search(ctx context.Context, req ReqContext, search string) ([]Group, error)

	// Gets a group by id.
	Get(ctx context.Context, req ReqContext, id ID) (Group, error)

	// Updates response for a group.
	Respond(ctx context.Context, req ReqContext, id ID, attendees uint8) error

	// ADMIN ACTIONS

	// Finds a group by name. This is a fuzzy search so should confirm before taking actions.
	FindOne(ctx context.Context, search string) (Group, error)

	// Updates a group's values. Only updates values that are present.
	Update(ctx context.Context, update GroupUpdate) error

	// Deletes a group.
	Delete(ctx context.Context, id ID) error

	// Imports a list of groups.
	Import(ctx context.Context, groups []GroupImport) error

	// Gets all groups for an export.
	Export(ctx context.Context) ([]Group, error)
}

type GroupRespository interface {
	Create(ctx context.Context, group Group) error
	CreateMany(ctx context.Context, groups []Group) error
	Update(ctx context.Context, group Group) error
	Get(ctx context.Context, id ID) (Group, error)
	Delete(ctx context.Context, id ID) error

	GetAll(ctx context.Context) ([]Group, error)

	SearchByName(ctx context.Context, search string) ([]Group, error)
}
