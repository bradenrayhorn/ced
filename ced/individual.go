package ced

import "context"

type Individual struct {
	ID           ID
	GroupID      ID
	Name         Name
	Response     bool
	HasResponded bool
}

type IndividualRespository interface {
	Create(ctx context.Context, individual Individual) error
	Get(ctx context.Context, id ID) (Individual, error)
}
