package sqlite_test

import (
	"context"
	"strings"
	"testing"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/internal/testutils"
	"github.com/bradenrayhorn/ced/sqlite"
	"github.com/matryer/is"
)

func TestIndividual(t *testing.T) {
	is := is.New(t)

	group := ced.Group{ID: ced.NewID()}
	group2 := ced.Group{ID: ced.NewID()}
	group3 := ced.Group{ID: ced.NewID()}
	var individualRepository ced.IndividualRespository
	var groupRepository ced.GroupRespository

	setup := func() func() {
		pool, stop := testutils.StartPool(t)
		groupRepository = sqlite.NewGroupRepository(pool)
		individualRepository = sqlite.NewIndividualRepository(pool)

		is.NoErr(groupRepository.Create(context.Background(), group))
		is.NoErr(groupRepository.Create(context.Background(), group2))
		is.NoErr(groupRepository.Create(context.Background(), group3))

		return stop
	}

	t.Run("can create and get", func(t *testing.T) {
		defer setup()()

		individual := ced.Individual{
			ID:           ced.NewID(),
			GroupID:      group.ID,
			Name:         ced.Name("Harry"),
			Response:     true,
			HasResponded: true,
		}
		err := individualRepository.Create(context.Background(), individual)
		is.NoErr(err)

		res, err := individualRepository.Get(context.Background(), individual.ID)
		is.NoErr(err)
		is.Equal(res, individual)
	})

	t.Run("cannot create if group does not exist", func(t *testing.T) {
		defer setup()()

		individual := ced.Individual{
			ID:           ced.NewID(),
			GroupID:      ced.NewID(),
			Name:         ced.Name("Harry"),
			Response:     true,
			HasResponded: true,
		}
		err := individualRepository.Create(context.Background(), individual)
		is.True(strings.Contains(err.Error(), "FOREIGN KEY constraint failed"))
	})

	t.Run("cannot create duplicates", func(t *testing.T) {
		defer setup()()

		individual := ced.Individual{
			ID:           ced.NewID(),
			GroupID:      group.ID,
			Name:         ced.Name("Harry"),
			Response:     true,
			HasResponded: true,
		}
		err := individualRepository.Create(context.Background(), individual)
		is.NoErr(err)

		err = individualRepository.Create(context.Background(), individual)
		is.True(err != nil)
	})

	t.Run("gives error when individual does not exist", func(t *testing.T) {
		defer setup()()

		_, err := individualRepository.Get(context.Background(), ced.NewID())
		testutils.IsCode(t, err, ced.ENOTFOUND)
	})

	t.Run("can update", func(t *testing.T) {
		defer setup()()

		individual := ced.Individual{
			ID:           ced.NewID(),
			GroupID:      group.ID,
			Name:         ced.Name("Harry"),
			Response:     true,
			HasResponded: true,
		}
		err := individualRepository.Create(context.Background(), individual)
		is.NoErr(err)

		res, err := individualRepository.Get(context.Background(), individual.ID)
		is.NoErr(err)
		is.Equal(res, individual)

		individual.GroupID = group2.ID
		individual.Name = ced.Name("Lacy")
		individual.Response = false
		individual.HasResponded = false

		err = individualRepository.Update(context.Background(), individual)
		is.NoErr(err)

		res, err = individualRepository.Get(context.Background(), individual.ID)
		is.NoErr(err)
		is.Equal(res, individual)
	})

	t.Run("can find by group", func(t *testing.T) {
		defer setup()()

		individual1 := ced.Individual{
			ID:      ced.NewID(),
			GroupID: group.ID,
			Name:    ced.Name("Harry"),
		}
		err := individualRepository.Create(context.Background(), individual1)
		is.NoErr(err)

		individual2 := ced.Individual{
			ID:      ced.NewID(),
			GroupID: group2.ID,
			Name:    ced.Name("Lacy"),
		}
		err = individualRepository.Create(context.Background(), individual2)
		is.NoErr(err)

		individual3 := ced.Individual{
			ID:      ced.NewID(),
			GroupID: group2.ID,
			Name:    ced.Name("George"),
		}
		err = individualRepository.Create(context.Background(), individual3)
		is.NoErr(err)

		res, err := individualRepository.GetByGroup(context.Background(), group2.ID)
		is.NoErr(err)
		is.Equal(
			testutils.SortSlice(res, testutils.CompareIndividuals),
			testutils.SortSlice([]ced.Individual{
				individual2,
				individual3,
			}, testutils.CompareIndividuals),
		)
	})

	t.Run("can search by name", func(t *testing.T) {
		defer setup()()

		individual1 := ced.Individual{
			ID:      ced.NewID(),
			GroupID: group.ID,
			Name:    ced.Name("Harry"),
		}
		err := individualRepository.Create(context.Background(), individual1)
		is.NoErr(err)

		individual2 := ced.Individual{
			ID:      ced.NewID(),
			GroupID: group2.ID,
			Name:    ced.Name("Lacy"),
		}
		err = individualRepository.Create(context.Background(), individual2)
		is.NoErr(err)

		individual3 := ced.Individual{
			ID:      ced.NewID(),
			GroupID: group2.ID,
			Name:    ced.Name("George"),
		}
		err = individualRepository.Create(context.Background(), individual3)
		is.NoErr(err)

		individual4 := ced.Individual{
			ID:      ced.NewID(),
			GroupID: group3.ID,
			Name:    ced.Name("Lacy"),
		}
		err = individualRepository.Create(context.Background(), individual4)
		is.NoErr(err)

		res, err := individualRepository.SearchByName(context.Background(), "lAcy ")
		is.NoErr(err)
		for k, v := range res {
			res[k] = testutils.SortSlice(v, testutils.CompareIndividuals)
		}

		is.Equal(
			res,
			map[ced.ID][]ced.Individual{
				individual3.GroupID: testutils.SortSlice([]ced.Individual{individual3, individual2}, testutils.CompareIndividuals),
				individual4.GroupID: {individual4},
			},
		)
	})
}
