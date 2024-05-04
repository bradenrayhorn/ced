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

func (c *groupContract) Search(ctx context.Context, req ced.ReqContext, search string) ([]ced.Group, error) {
	res, err := c.groupRepository.SearchByName(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("search groups %s: %w", search, err)
	}

	slog.Info("searched groups",
		"search", search,
		"results", len(res),
		"ip", req.ConnectingIP,
	)

	return res, nil
}

func (c *groupContract) Get(ctx context.Context, req ced.ReqContext, id ced.ID) (ced.Group, error) {
	res, err := c.groupRepository.Get(ctx, id)
	if err != nil {
		if !errors.Is(err, ced.ErrorNotFound) {
			return res, fmt.Errorf("get group %s: %w", id, err)
		}

		slog.Info("get group not found",
			"id", id,
			"ip", req.ConnectingIP,
		)

		return res, err
	}
	return res, nil
}

func (c *groupContract) Respond(ctx context.Context, req ced.ReqContext, id ced.ID, attendees uint8) error {
	request := struct {
		Id        ced.ID
		Attendees uint8
		Req       ced.ReqContext
	}{id, attendees, req}

	group, err := c.groupRepository.Get(ctx, id)
	if err != nil {
		if !errors.Is(err, ced.ErrorNotFound) {
			return fmt.Errorf("get group to respond %+v: %w", request, err)
		}
		return err
	}

	group.Attendees = attendees
	group.HasResponded = true

	if err = group.Validate(); err != nil {
		return err
	}

	slog.Info("group updated",
		"attendees", attendees,
		"id", id,
		"name", group.Name,
		"ip", req.ConnectingIP,
	)

	err = c.groupRepository.Update(ctx, group)
	if err != nil {
		return fmt.Errorf("respond %+v to group %+v: %w", request, group, err)
	}
	return nil
}

func (c *groupContract) FindOne(ctx context.Context, search string) (ced.Group, error) {
	res, err := c.groupRepository.SearchByName(ctx, search)
	if err != nil {
		return ced.Group{}, fmt.Errorf("find one group %s: %w", search, err)
	}

	if len(res) == 0 {
		return ced.Group{}, fmt.Errorf("find one group: no results")
	}

	if len(res) > 1 {
		return ced.Group{}, fmt.Errorf("find one group: too many results, try narrowing the search")
	}

	return res[0], nil
}

func (c *groupContract) Update(ctx context.Context, update ced.GroupUpdate) error {
	group, err := c.groupRepository.Get(ctx, update.ID)
	if err != nil {
		return fmt.Errorf("update group: %w", err)
	}

	if update.Attendees != nil {
		group.Attendees = *update.Attendees
		group.HasResponded = true
	}
	if update.MaxAttendees != nil {
		group.MaxAttendees = *update.MaxAttendees
	}
	if update.Name != nil {
		group.Name = *update.Name
	}
	if update.SearchHints != nil {
		group.SearchHints = *update.SearchHints
	}

	if err := group.Validate(); err != nil {
		var cedError ced.Error
		if errors.As(err, &cedError) {
			_, msg := cedError.CedError()
			return fmt.Errorf("update group: %s", msg)
		}
		return err
	}

	if err := c.groupRepository.Update(ctx, group); err != nil {
		return fmt.Errorf("update group: %w", err)
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

func (c *groupContract) Export(ctx context.Context) ([]ced.Group, error) {
	return c.groupRepository.GetAll(ctx)
}
