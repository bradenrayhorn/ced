package contract

import (
	"context"

	"github.com/bradenrayhorn/ced/ced"
)

var _ ced.IndividualContract = (*individualContract)(nil)

type individualContract struct {
	individualRepository ced.IndividualRespository
}

func NewIndividualContract(
	individualRepository ced.IndividualRespository,
) *individualContract {
	return &individualContract{individualRepository}
}

func (c *individualContract) SearchByName(ctx context.Context, search string) (map[ced.ID][]ced.Individual, error) {
	return c.individualRepository.SearchByName(ctx, search)
}

func (c *individualContract) GetInGroup(ctx context.Context, groupID ced.ID) ([]ced.Individual, error) {
	return c.individualRepository.GetByGroup(ctx, groupID)
}

func (c *individualContract) SetResponse(ctx context.Context, individualID ced.ID, response bool) error {
	individual, err := c.individualRepository.Get(ctx, individualID)
	if err != nil {
		return err
	}

	individual.Response = response
	individual.HasResponded = true

	return c.individualRepository.Update(ctx, individual)
}
