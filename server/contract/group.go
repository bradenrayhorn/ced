package contract

import (
	"context"
	"fmt"

	"github.com/bradenrayhorn/ced/server/ced"
)

var _ ced.GroupContract = (*groupContract)(nil)

type groupContract struct {
	groupRepository ced.GroupRespository
}

func NewGroupContract(
	groupRepository ced.GroupRespository,
) *groupContract {
	return &groupContract{groupRepository}
}

func (c *groupContract) Create(ctx context.Context, name ced.Name, maxAttendees uint8) (ced.Group, error) {
	if err := ced.ValidateFields(ced.Field("Name", name)); err != nil {
		return ced.Group{}, err
	}

	group := ced.Group{
		ID:           ced.NewID(),
		Name:         name,
		MaxAttendees: maxAttendees,
		Attendees:    0,
		HasResponded: false,
	}

	if err := c.groupRepository.Create(ctx, group); err != nil {
		return ced.Group{}, err
	}

	return group, nil
}

func (c *groupContract) Search(ctx context.Context, search string) ([]ced.Group, error) {
	return c.groupRepository.SearchByName(ctx, search)
}

func (c *groupContract) Get(ctx context.Context, id ced.ID) (ced.Group, error) {
	return c.groupRepository.Get(ctx, id)
}

func (c *groupContract) Respond(ctx context.Context, id ced.ID, attendees uint8) error {
	group, err := c.groupRepository.Get(ctx, id)
	if err != nil {
		return err
	}

	if attendees > group.MaxAttendees {
		return ced.NewError(ced.EINVALID, fmt.Sprintf("group can have at most %d attendees", group.MaxAttendees))
	}

	group.Attendees = attendees
	group.HasResponded = true

	return c.groupRepository.Update(ctx, group)
}
