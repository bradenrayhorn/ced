package ced

import "context"

type Group struct {
	ID           ID
	Name         Name
	Attendees    uint8
	MaxAttendees uint8
	HasResponded bool
}

type GroupContract interface {
	// Creates a new group.
	Create(ctx context.Context, name Name, maxAttendees uint8) (Group, error)

	// Searches for a group using the search string.
	Search(ctx context.Context, search string) ([]Group, error)

	// Gets a group by id.
	Get(ctx context.Context, id ID) (Group, error)

	// Updates response for a group.
	Respond(ctx context.Context, id ID, attendees uint8, connectingIP string) error
}

type GroupRespository interface {
	Create(ctx context.Context, group Group) error
	Update(ctx context.Context, group Group) error
	Get(ctx context.Context, id ID) (Group, error)

	SearchByName(ctx context.Context, search string) ([]Group, error)
}
