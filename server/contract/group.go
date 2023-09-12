package contract

import (
	"context"
	"fmt"
	"log/slog"

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

func (c *groupContract) Create(ctx context.Context, name ced.Name, maxAttendees uint8, searchHints string) (ced.Group, error) {
	group, err := ced.NewGroup(name, maxAttendees, searchHints)
	if err != nil {
		return ced.Group{}, err
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

func (c *groupContract) Respond(ctx context.Context, id ced.ID, attendees uint8, connectingIP string) error {
	group, err := c.groupRepository.Get(ctx, id)
	if err != nil {
		return err
	}

	if attendees > group.MaxAttendees {
		return ced.NewError(ced.EINVALID, fmt.Sprintf("group can have at most %d attendees", group.MaxAttendees))
	}

	group.Attendees = attendees
	group.HasResponded = true

	slog.Info("group updated",
		"attendees", attendees,
		"id", id,
		"name", group.Name,
		"ip", connectingIP,
	)

	return c.groupRepository.Update(ctx, group)
}

func (c *groupContract) Import(ctx context.Context, records []ced.GroupImport) error {
	groups := make([]ced.Group, len(records))
	for i, record := range records {
		group, err := ced.NewGroup(record.Name, record.MaxAttendees, record.SearchHints)
		if err != nil {
			return fmt.Errorf("failed to import at record %d: %w", i+1, err)
		}

		groups[i] = group
	}

	return c.groupRepository.CreateMany(ctx, groups)
}
