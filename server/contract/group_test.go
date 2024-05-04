package contract_test

import (
	"context"
	"fmt"
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
		SearchHints:  "George Hoover, Tom Thumb",
	}
	group2 := ced.Group{
		ID:           ced.NewID(),
		Name:         "Elizabeth Hoover and family",
		MaxAttendees: 2,
		Attendees:    0,
		HasResponded: false,
		SearchHints:  "Elizabeth Hoover, Tom Thumb",
	}
	reqCtx := ced.ReqContext{}

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

		res, err := groupContract.Search(context.Background(), reqCtx, "George Hover")
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
		t.Run("can update response", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			err := groupContract.Respond(context.Background(), reqCtx, group1.ID, 5)
			is.NoErr(err)

			res, err := groupRepository.Get(context.Background(), group1.ID)
			is.NoErr(err)
			is.Equal(res.HasResponded, true)
			is.Equal(res.Attendees, uint8(5))
		})

		t.Run("can update response to not going", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			err := groupContract.Respond(context.Background(), reqCtx, group2.ID, 0)
			is.NoErr(err)

			res, err := groupRepository.Get(context.Background(), group2.ID)
			is.NoErr(err)
			is.Equal(res.HasResponded, true)
			is.Equal(res.Attendees, uint8(0))
		})

		t.Run("cannot update invalid group", func(t *testing.T) {
			defer setup(t)()

			id := ced.NewID()
			err := groupContract.Respond(context.Background(), reqCtx, id, 1)
			testutils.IsCodeAndError(t, err, ced.ENOTFOUND, fmt.Sprintf("ced.Group [%s] not found", id))
		})

		t.Run("cannot update to more attendees than allowed", func(t *testing.T) {
			defer setup(t)()

			err := groupContract.Respond(context.Background(), reqCtx, group1.ID, 6)
			testutils.IsCodeAndError(t, err, ced.EINVALID, "group can have at most 5 attendees, has 6")
		})
	})

	t.Run("get group", func(t *testing.T) {

		t.Run("can get group", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			res, err := groupContract.Get(context.Background(), reqCtx, group1.ID)
			is.NoErr(err)
			is.Equal(res, group1)
		})

		t.Run("can get 404", func(t *testing.T) {
			defer setup(t)()

			_, err := groupContract.Get(context.Background(), reqCtx, ced.NewID())
			testutils.IsCode(t, err, ced.ENOTFOUND)
		})
	})

	t.Run("find one", func(t *testing.T) {

		t.Run("can find one group", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			res, err := groupContract.FindOne(context.Background(), "Goerge Hoover and family")
			is.NoErr(err)
			is.Equal(res, group1)
		})

		t.Run("handles if nothing found", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			_, err := groupContract.FindOne(context.Background(), "nonsense")
			is.Equal("find one group: no results", err.Error())
		})

		t.Run("handles too many results", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			_, err := groupContract.FindOne(context.Background(), "Tom Thumb")
			is.Equal("find one group: too many results, try narrowing the search", err.Error())
		})
	})

	t.Run("update group", func(t *testing.T) {

		t.Run("handles group not found", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			update := ced.GroupUpdate{ID: ced.NewID()}

			err := groupContract.Update(context.Background(), update)
			is.Equal("update group: Not found", err.Error())
		})

		t.Run("validates name", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			update := ced.GroupUpdate{
				ID:   group1.ID,
				Name: pointerize(ced.Name("")),
			}

			err := groupContract.Update(context.Background(), update)
			is.Equal("update group: Name is required.", err.Error())
		})

		t.Run("validates attendees update", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			update := ced.GroupUpdate{
				ID:        group1.ID,
				Attendees: pointerize(uint8(6)),
			}

			err := groupContract.Update(context.Background(), update)
			is.Equal("update group: group can have at most 5 attendees, has 6", err.Error())
		})

		t.Run("validates max attendees update", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			update := ced.GroupUpdate{
				ID:           group1.ID,
				MaxAttendees: pointerize(uint8(0)),
			}

			err := groupContract.Update(context.Background(), update)
			is.Equal("update group: group can have at most 0 attendees, has 4", err.Error())
		})

		t.Run("validates attendees and max attendees update", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			update := ced.GroupUpdate{
				ID:           group1.ID,
				Attendees:    pointerize(uint8(2)),
				MaxAttendees: pointerize(uint8(1)),
			}

			err := groupContract.Update(context.Background(), update)
			is.Equal("update group: group can have at most 1 attendees, has 2", err.Error())
		})

		t.Run("can update name", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			update := ced.GroupUpdate{
				ID:   group1.ID,
				Name: pointerize(ced.Name("The Muffin Man")),
			}

			err := groupContract.Update(context.Background(), update)
			is.NoErr(err)

			res, err := groupRepository.Get(context.Background(), group1.ID)
			is.NoErr(err)

			is.Equal(ced.Name("The Muffin Man"), res.Name)
			is.Equal(group1.Attendees, res.Attendees)
			is.Equal(group1.MaxAttendees, res.MaxAttendees)
			is.Equal(group1.HasResponded, res.HasResponded)
			is.Equal(group1.SearchHints, res.SearchHints)
		})

		t.Run("can update attendees", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			update := ced.GroupUpdate{
				ID:        group2.ID,
				Attendees: pointerize(uint8(1)),
			}

			err := groupContract.Update(context.Background(), update)
			is.NoErr(err)

			res, err := groupRepository.Get(context.Background(), group2.ID)
			is.NoErr(err)

			is.Equal(group2.Name, res.Name)
			is.Equal(uint8(1), res.Attendees)
			is.Equal(group2.MaxAttendees, res.MaxAttendees)
			is.Equal(true, res.HasResponded)
			is.Equal(group2.SearchHints, res.SearchHints)
		})

		t.Run("can update max attendees", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			update := ced.GroupUpdate{
				ID:           group2.ID,
				MaxAttendees: pointerize(uint8(10)),
			}

			err := groupContract.Update(context.Background(), update)
			is.NoErr(err)

			res, err := groupRepository.Get(context.Background(), group2.ID)
			is.NoErr(err)

			is.Equal(group2.Name, res.Name)
			is.Equal(group2.Attendees, res.Attendees)
			is.Equal(uint8(10), res.MaxAttendees)
			is.Equal(false, res.HasResponded)
			is.Equal(group2.SearchHints, res.SearchHints)
		})

		t.Run("can update search hints", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			update := ced.GroupUpdate{
				ID:          group2.ID,
				SearchHints: pointerize("here is a hint"),
			}

			err := groupContract.Update(context.Background(), update)
			is.NoErr(err)

			res, err := groupRepository.Get(context.Background(), group2.ID)
			is.NoErr(err)

			is.Equal(group2.Name, res.Name)
			is.Equal(group2.Attendees, res.Attendees)
			is.Equal(group2.MaxAttendees, res.MaxAttendees)
			is.Equal(group2.HasResponded, res.HasResponded)
			is.Equal("here is a hint", res.SearchHints)
		})
	})

	t.Run("import", func(t *testing.T) {

		t.Run("can import", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			groupImport := ced.GroupImport{
				Name:         ced.Name("Bob"),
				MaxAttendees: 2,
				SearchHints:  "Bob Lob",
			}

			err := groupContract.Import(context.Background(), []ced.GroupImport{groupImport})
			is.NoErr(err)

			res, err := groupContract.Search(context.Background(), reqCtx, "Bob Lob")
			is.NoErr(err)
			is.Equal(len(res), 1)
			is.Equal(res[0].Name, ced.Name("Bob"))
			is.Equal(res[0].MaxAttendees, uint8(2))
			is.Equal(res[0].SearchHints, "Bob Lob")
		})

		t.Run("when fails adds record number", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			groupImport := ced.GroupImport{
				Name:         ced.Name(""),
				MaxAttendees: 2,
				SearchHints:  "Bob Lob",
			}

			err := groupContract.Import(context.Background(), []ced.GroupImport{groupImport})
			is.True(err != nil)
			is.Equal(err.Error(), "failed to import at record 1: Invalid data provided")
		})
	})

	t.Run("export", func(t *testing.T) {
		t.Run("can export", func(t *testing.T) {
			defer setup(t)()
			is := is.New(t)

			res, err := groupContract.Export(context.Background())
			is.NoErr(err)

			is.Equal(2, len(res))

			var resGroup1 ced.Group
			var resGroup2 ced.Group
			for _, group := range res {
				if group.ID == group1.ID {
					resGroup1 = group
				}
				if group.ID == group2.ID {
					resGroup2 = group
				}
			}

			is.Equal(group1, resGroup1)
			is.Equal(group2, resGroup2)
		})
	})
}

func pointerize[T any](t T) *T {
	return &t
}
