package ced

import "context"

type Individual struct {
	ID           ID
	GroupID      ID
	Name         Name
	Response     bool
	HasResponded bool
}

type IndividualContract interface {
	// Finds all individuals matching the search text.
	SearchByName(ctx context.Context, search string) (map[ID][]Individual, error)

	// Lists all individuals in a group.
	GetInGroup(ctx context.Context, groupID ID) ([]Individual, error)

	// Update an individual's response.
	Respond(ctx context.Context, individualID ID, response bool) error
}

type IndividualRespository interface {
	Create(ctx context.Context, individual Individual) error
	Get(ctx context.Context, id ID) (Individual, error)
	Update(ctx context.Context, individual Individual) error

	GetByGroup(ctx context.Context, groupID ID) ([]Individual, error)
	SearchByName(ctx context.Context, search string) (map[ID][]Individual, error)
}
