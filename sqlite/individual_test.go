package sqlite_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/internal/testutils"
	"github.com/bradenrayhorn/ced/sqlite"
	"github.com/matryer/is"
)

func TestIndividual(t *testing.T) {
	is := is.New(t)

	group := ced.Group{ID: ced.NewID()}
	var individualRepository ced.IndividualRespository
	var groupRepository ced.GroupRespository

	setup := func() func() {
		pool, stop := testutils.StartPool(t)
		groupRepository = sqlite.NewGroupRepository(pool)
		individualRepository = sqlite.NewIndividualRepository(pool)

		is.NoErr(groupRepository.Create(context.Background(), group))

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
}
