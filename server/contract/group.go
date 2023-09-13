package contract

import (
	"context"
	"errors"
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
		return ced.Group{}, fmt.Errorf("create group %+v: %w", group, err)
	}

	return group, nil
}

func (c *groupContract) Search(ctx context.Context, search string) ([]ced.Group, error) {
	res, err := c.groupRepository.SearchByName(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("search groups %s: %w", search, err)
	}
	return res, nil
}

func (c *groupContract) Get(ctx context.Context, id ced.ID) (ced.Group, error) {
	res, err := c.groupRepository.Get(ctx, id)
	if err != nil {
		if !errors.Is(err, ced.ErrorNotFound) {
			return res, fmt.Errorf("get group %s: %w", id, err)
		}
		return res, err
	}
	return res, nil
}

func (c *groupContract) Respond(ctx context.Context, id ced.ID, attendees uint8, connectingIP string) error {
	req := struct {
		id           ced.ID
		attendees    uint8
		connectingIP string
	}{id, attendees, connectingIP}

	group, err := c.groupRepository.Get(ctx, id)
	if err != nil {
		if !errors.Is(err, ced.ErrorNotFound) {
			return fmt.Errorf("get group to respond %+v: %w", req, err)
		}
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

	err = c.groupRepository.Update(ctx, group)
	if err != nil {
		return fmt.Errorf("respond %+v to group %+v: %w", req, group, err)
	}
	return nil
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
