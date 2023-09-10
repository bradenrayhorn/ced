package contract_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/ced/server/ced"
	"github.com/bradenrayhorn/ced/server/contract"
	"github.com/bradenrayhorn/ced/server/internal/testutils"
	"github.com/bradenrayhorn/ced/server/sqlite"
	"github.com/matryer/is"
)

func TestGroup(t *testing.T) {
	group1 := ced.Group{
		ID:           ced.NewID(),
		Name:         "George Hoover and family",
		MaxAttendees: 5,
		Attendees:    4,
		HasResponded: true,
		SearchHints:  "George Hoover",
	}
	group2 := ced.Group{
		ID:           ced.NewID(),
		Name:         "Elizabeth Hoover and family",
		MaxAttendees: 2,
		Attendees:    0,
		HasResponded: false,
		SearchHints:  "Elizabeth Hoover",
	}

	var groupRepository ced.GroupRespository
	var groupContract ced.GroupContract

	setup := func(t *testing.T) func() {
		is := is.New(t)
		pool, stop := testutils.StartPool(t)
		groupRepository = sqlite.NewGroupRepository(pool)

		groupContract = contract.NewGroupContract(groupRepository)

		is.NoErr(groupRepository.Create(context.Background(), group1))
		is.NoErr(groupRepository.Create(context.Background(), group2))

		return stop
	}

	t.Run("can search by name", func(t *testing.T) {
		defer setup(t)()
		is := is.New(t)

		res, err := groupContract.Search(context.Background(), "George Hover")
		is.NoErr(err)

		is.Equal(res, []ced.Group{group1})
	})

	t.Run("create group", func(t *testing.T) {

		t.Run("can create", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			res, err := groupContract.Create(context.Background(), ced.Name("Tom"), 1, "Thomas")
			is.NoErr(err)
			is.True(!res.ID.Empty())
			is.Equal(res.Name, ced.Name("Tom"))
			is.Equal(res.Attendees, uint8(0))
			is.Equal(res.MaxAttendees, uint8(1))
			is.Equal(res.HasResponded, false)
			is.Equal(res.SearchHints, "Thomas")
		})

		t.Run("name must be valid", func(t *testing.T) {
			defer setup(t)()

			_, err := groupContract.Create(context.Background(), ced.Name(""), 1, "")
			testutils.IsCodeAndError(t, err, ced.EINVALID, "Name is required.")
		})
	})

	t.Run("update group", func(t *testing.T) {
		ip := "127.0.0.1"

		t.Run("can update response", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			err := groupContract.Respond(context.Background(), group1.ID, 5, ip)
			is.NoErr(err)

			res, err := groupRepository.Get(context.Background(), group1.ID)
			is.NoErr(err)
			is.Equal(res.HasResponded, true)
			is.Equal(res.Attendees, uint8(5))
		})

		t.Run("can update response to not going", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			err := groupContract.Respond(context.Background(), group2.ID, 0, ip)
			is.NoErr(err)

			res, err := groupRepository.Get(context.Background(), group2.ID)
			is.NoErr(err)
			is.Equal(res.HasResponded, true)
			is.Equal(res.Attendees, uint8(0))
		})

		t.Run("cannot update invalid group", func(t *testing.T) {
			defer setup(t)()

			err := groupContract.Respond(context.Background(), ced.NewID(), 1, ip)
			testutils.IsCodeAndError(t, err, ced.ENOTFOUND, "ced.Group not found.")
		})

		t.Run("cannot update to more attendees than allowed", func(t *testing.T) {
			defer setup(t)()

			err := groupContract.Respond(context.Background(), group1.ID, 6, ip)
			testutils.IsCodeAndError(t, err, ced.EINVALID, "group can have at most 5 attendees")
		})
	})

	t.Run("get group", func(t *testing.T) {

		t.Run("can get group", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			res, err := groupContract.Get(context.Background(), group1.ID)
			is.NoErr(err)
			is.Equal(res, group1)
		})
	})
}
