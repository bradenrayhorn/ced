package ced

import "context"

type Group struct {
	ID ID
}

type GroupContract interface {
	// Creates a new group. All individuals are created and added to the group.
	Create(ctx context.Context, individuals []Name) (Group, error)
}

type GroupRespository interface {
	Create(ctx context.Context, group Group) error
	Get(ctx context.Context, id ID) (Group, error)
}
