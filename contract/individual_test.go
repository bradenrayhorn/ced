package contract_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/contract"
	"github.com/bradenrayhorn/ced/internal/testutils"
	"github.com/bradenrayhorn/ced/sqlite"
	"github.com/matryer/is"
)

func TestIndividual(t *testing.T) {
	group1 := ced.Group{ID: ced.NewID()}
	group2 := ced.Group{ID: ced.NewID()}
	individual1 := ced.Individual{ID: ced.NewID(), GroupID: group1.ID, Name: "Harry", Response: false, HasResponded: false}
	individual2 := ced.Individual{ID: ced.NewID(), GroupID: group2.ID, Name: "Lacy", Response: false, HasResponded: true}
	individual3 := ced.Individual{ID: ced.NewID(), GroupID: group2.ID, Name: "George", Response: true, HasResponded: false}

	var individualRepository ced.IndividualRespository
	var individualContract ced.IndividualContract

	setup := func(t *testing.T) func() {
		is := is.New(t)
		pool, stop := testutils.StartPool(t)
		groupRepository := sqlite.NewGroupRepository(pool)
		individualRepository = sqlite.NewIndividualRepository(pool)

		individualContract = contract.NewIndividualContract(individualRepository)

		is.NoErr(groupRepository.Create(context.Background(), group1))
		is.NoErr(groupRepository.Create(context.Background(), group2))

		is.NoErr(individualRepository.Create(context.Background(), individual1))
		is.NoErr(individualRepository.Create(context.Background(), individual2))
		is.NoErr(individualRepository.Create(context.Background(), individual3))

		return stop
	}

	t.Run("can search by name", func(t *testing.T) {
		defer setup(t)()
		is := is.New(t)

		res, err := individualContract.SearchByName(context.Background(), "Lacy")
		is.NoErr(err)
		for k, v := range res {
			res[k] = testutils.SortSlice(v, testutils.CompareIndividuals)
		}

		is.Equal(res, map[ced.ID][]ced.Individual{
			group2.ID: testutils.SortSlice([]ced.Individual{individual2, individual3}, testutils.CompareIndividuals),
		})
	})

	t.Run("can get by group", func(t *testing.T) {
		defer setup(t)()
		is := is.New(t)

		res, err := individualContract.GetInGroup(context.Background(), group1.ID)
		is.NoErr(err)
		is.Equal(res, []ced.Individual{individual1})
	})

	t.Run("can update response to true", func(t *testing.T) {
		defer setup(t)()
		is := is.New(t)

		err := individualContract.SetResponse(context.Background(), individual1.ID, true)
		is.NoErr(err)

		res, err := individualRepository.Get(context.Background(), individual1.ID)
		is.NoErr(err)
		is.Equal(res.Response, true)
		is.Equal(res.HasResponded, true)
	})

	t.Run("can update response to false", func(t *testing.T) {
		defer setup(t)()
		is := is.New(t)

		err := individualContract.SetResponse(context.Background(), individual2.ID, false)
		is.NoErr(err)

		res, err := individualRepository.Get(context.Background(), individual2.ID)
		is.NoErr(err)
		is.Equal(res.Response, false)
		is.Equal(res.HasResponded, true)
	})
}
