package sqlite_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/ced/server/ced"
	"github.com/bradenrayhorn/ced/server/internal/testutils"
	"github.com/bradenrayhorn/ced/server/sqlite"
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
			SearchHints:  "George Hoover, Eleanor Hoover",
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

	t.Run("create many operates within a transaction", func(t *testing.T) {
		is := is.New(t)
		defer setup()()

		// group3 fails to create because it shares the same ID as group1.
		// this should rollback the entire create operation.
		group1 := ced.Group{ID: ced.NewID()}
		group2 := ced.Group{ID: ced.NewID()}
		group3 := ced.Group{ID: group1.ID}
		err := groupRepository.CreateMany(context.Background(), []ced.Group{group1, group2, group3})
		is.True(err != nil)

		_, err = groupRepository.Get(context.Background(), group1.ID)
		testutils.IsCode(t, err, ced.ENOTFOUND)

		_, err = groupRepository.Get(context.Background(), group2.ID)
		testutils.IsCode(t, err, ced.ENOTFOUND)

		_, err = groupRepository.Get(context.Background(), group3.ID)
		testutils.IsCode(t, err, ced.ENOTFOUND)
	})

	t.Run("gives error when group does not exist", func(t *testing.T) {
		defer setup()()

		_, err := groupRepository.Get(context.Background(), ced.NewID())
		testutils.IsCode(t, err, ced.ENOTFOUND)
	})

	t.Run("can delete group", func(t *testing.T) {
		is := is.New(t)
		defer setup()()

		group := ced.Group{
			ID:   ced.NewID(),
			Name: "George Hoover and family",
		}
		err := groupRepository.Create(context.Background(), group)
		is.NoErr(err)

		err = groupRepository.Delete(context.Background(), group.ID)
		is.NoErr(err)

		_, err = groupRepository.Get(context.Background(), group.ID)
		testutils.IsCode(t, err, ced.ENOTFOUND)
	})

	t.Run("can get all groups", func(t *testing.T) {
		is := is.New(t)
		defer setup()()

		group := ced.Group{
			ID:           ced.NewID(),
			Name:         "George Hoover and family",
			Attendees:    3,
			MaxAttendees: 5,
			HasResponded: true,
		}

		is.NoErr(groupRepository.Create(context.Background(), group))

		res, err := groupRepository.GetAll(context.Background())
		is.NoErr(err)
		is.Equal([]ced.Group{group}, res)
	})

	t.Run("search", func(t *testing.T) {
		is := is.New(t)
		defer setup()()

		group1 := ced.Group{
			ID:          ced.NewID(),
			Name:        " George Hoover and family",
			SearchHints: "George Hoover, Kelly Hoover, Elizabeth Frank",
		}
		group2 := ced.Group{
			ID:          ced.NewID(),
			Name:        "Elizabeth George",
			SearchHints: "Elizabeth George, George Manthin, Poppy Seed, Poppy Seed",
		}
		group3 := ced.Group{
			ID:          ced.NewID(),
			Name:        "Geoff Kee and me",
			SearchHints: "Geoff Kee, Kelly Free",
		}
		group4 := ced.Group{
			ID:          ced.NewID(),
			Name:        "Geoff Tree & company",
			SearchHints: "Geoff Tree, Kelly Free",
		}
		is.NoErr(groupRepository.Create(context.Background(), group1))
		is.NoErr(groupRepository.Create(context.Background(), group2))
		is.NoErr(groupRepository.Create(context.Background(), group3))
		is.NoErr(groupRepository.Create(context.Background(), group4))

		var tests = []struct {
			name           string
			search         string
			expectedGroups []ced.Group
		}{
			{"finds unique name",
				"George Hoover", []ced.Group{group1}},
			{"cannot find anything with duplicate last name",
				"Hoover", []ced.Group{}},
			{"can correct a typo",
				"Geoge Hoover", []ced.Group{group1}},
			{"can correct a 3-character typo",
				"Geoge Hovet", []ced.Group{group1}},
			{"cannot correct a 4-character typo",
				"Geofe Hovet", []ced.Group{}},
			{"does not include same group twice",
				"Poppy Seed", []ced.Group{group2}},
			{"may include two groups",
				"Geoff Keee", []ced.Group{group3, group4}},
			{"stops matching if finds exact match",
				"Geoff Kee", []ced.Group{group3}},
			{"searches on group name",
				"Geoff Kee and me", []ced.Group{group3}},
			{"can include multiple exact matches",
				"Kelly Free", []ced.Group{group3, group4}},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				is := is.New(t)

				res, err := groupRepository.SearchByName(context.Background(), test.search)
				is.NoErr(err)
				is.Equal(res, test.expectedGroups)
			})
		}
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
			SearchHints:  "George Hoover, Eleanor Hoover",
		}
		is.NoErr(groupRepository.Create(context.Background(), group))

		group.Name = "Elizabeth Hoover and family"
		group.MaxAttendees = 4
		group.Attendees = 2
		group.HasResponded = false
		group.SearchHints = "Elizabeth Hoover, Eleanor Hoover"

		is.NoErr(groupRepository.Update(context.Background(), group))

		res, err := groupRepository.Get(context.Background(), group.ID)
		is.NoErr(err)
		is.Equal(res, group)
	})
}
