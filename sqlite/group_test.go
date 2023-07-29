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
	is := is.New(t)

	var groupRepository ced.GroupRespository

	setup := func() func() {
		pool, stop := testutils.StartPool(t)
		groupRepository = sqlite.NewGroupRepository(pool)
		return stop
	}

	t.Run("can create and get", func(t *testing.T) {
		defer setup()()

		group := ced.Group{ID: ced.NewID()}
		err := groupRepository.Create(context.Background(), group)
		is.NoErr(err)

		res, err := groupRepository.Get(context.Background(), group.ID)
		is.NoErr(err)
		is.Equal(res, group)
	})

	t.Run("cannot create duplicates", func(t *testing.T) {
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
}
