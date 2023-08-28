package sqlite_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/internal/testutils"
	"github.com/bradenrayhorn/ced/sqlite"
	"github.com/matryer/is"
)

func TestGroup(t *testing.T) {
	var groupRepository ced.GroupRespository

	setup := func() func() {
		pool, stop := testutils.StartPool(t)
		groupRepository = sqlite.NewGroupRepository(pool)
		return stop
	}

	t.Run("can create and get", func(t *testing.T) {
		is := is.New(t)
		defer setup()()

		group := ced.Group{
			ID:           ced.NewID(),
			Name:         "George Hoover and family",
			MaxAttendees: 5,
			Attendees:    4,
			HasResponded: true,
		}
		err := groupRepository.Create(context.Background(), group)
		is.NoErr(err)

		res, err := groupRepository.Get(context.Background(), group.ID)
		is.NoErr(err)
		is.Equal(res, group)
	})

	t.Run("cannot create duplicates", func(t *testing.T) {
		is := is.New(t)
		defer setup()()

		group := ced.Group{ID: ced.NewID()}
		err := groupRepository.Create(context.Background(), group)
		is.NoErr(err)

		err = groupRepository.Create(context.Background(), group)
		is.True(err != nil)
	})

	t.Run("gives error when group does not exist", func(t *testing.T) {
		defer setup()()

		_, err := groupRepository.Get(context.Background(), ced.NewID())
		testutils.IsCode(t, err, ced.ENOTFOUND)
	})

	t.Run("search for group", func(t *testing.T) {
		is := is.New(t)
		defer setup()()

		group1 := ced.Group{
			ID:   ced.NewID(),
			Name: " George Hoover and family",
		}
		group2 := ced.Group{
			ID:   ced.NewID(),
			Name: "Elizabeth George",
		}
		group3 := ced.Group{
			ID:   ced.NewID(),
			Name: "Geoff Kee",
		}
		is.NoErr(groupRepository.Create(context.Background(), group1))
		is.NoErr(groupRepository.Create(context.Background(), group2))
		is.NoErr(groupRepository.Create(context.Background(), group3))

		res, err := groupRepository.SearchByName(context.Background(), " gEorGe")
		is.NoErr(err)
		is.Equal(
			testutils.SortSlice(res, testutils.CompareGroups),
			testutils.SortSlice([]ced.Group{group1, group2}, testutils.CompareGroups),
		)
	})

	t.Run("can update group", func(t *testing.T) {
		is := is.New(t)
		defer setup()()

		group := ced.Group{
			ID:           ced.NewID(),
			Name:         "George Hoover and family",
			MaxAttendees: 5,
			Attendees:    4,
			HasResponded: true,
		}
		is.NoErr(groupRepository.Create(context.Background(), group))

		group.Name = "Elizabeth Hoover and family"
		group.MaxAttendees = 4
		group.Attendees = 2
		group.HasResponded = false

		is.NoErr(groupRepository.Update(context.Background(), group))

		res, err := groupRepository.Get(context.Background(), group.ID)
		is.NoErr(err)
		is.Equal(res, group)
	})
}
