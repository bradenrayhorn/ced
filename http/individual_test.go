package http_test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/matryer/is"
)

func TestIndividualSearch(t *testing.T) {
	t.Run("param is required", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		r := test.DoRequest(t, "GET", "/api/v1/individuals/search", nil, http.StatusUnprocessableEntity)
		is.Equal(r, `{"error":"search is required","code":"unprocessable"}`)
	})

	t.Run("can search", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		individual := ced.Individual{
			ID:           ced.NewID(),
			GroupID:      ced.NewID(),
			Name:         "George",
			Response:     true,
			HasResponded: false,
		}
		test.individualContract.SearchByNameFunc.PushReturn(map[ced.ID][]ced.Individual{
			individual.GroupID: {individual},
		}, nil)

		r := test.DoRequest(t, "GET", "/api/v1/individuals/search?search=geo", nil, http.StatusOK)
		is.Equal(r, fmt.Sprintf(`{"data":{"%s":[{"id":"%s","name":"George","response":true,"has_responded":false}]}}`, individual.GroupID, individual.ID))

		params := test.individualContract.SearchByNameFunc.History()[0]
		is.Equal(params.Arg1, "geo")
	})

	t.Run("can handle search error", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		test.individualContract.SearchByNameFunc.PushReturn(nil, errors.New("query failed"))

		r := test.DoRequest(t, "GET", "/api/v1/individuals/search?search=geo", nil, http.StatusInternalServerError)
		is.Equal(r, `{"error":"Internal error","code":"internal"}`)
	})
}

func TestGetIndividualsInGroup(t *testing.T) {
	t.Run("param must be valid", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		r := test.DoRequest(t, "GET", "/api/v1/groups/--/individuals", nil, http.StatusUnprocessableEntity)
		is.Equal(r, `{"error":"invalid group id","code":"unprocessable"}`)
	})

	t.Run("can get in group", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		individual := ced.Individual{
			ID:           ced.NewID(),
			GroupID:      ced.NewID(),
			Name:         "George",
			Response:     true,
			HasResponded: false,
		}
		test.individualContract.GetInGroupFunc.PushReturn([]ced.Individual{individual}, nil)

		r := test.DoRequest(t, "GET", fmt.Sprintf("/api/v1/groups/%s/individuals", individual.GroupID), nil, http.StatusOK)
		is.Equal(r, fmt.Sprintf(`{"data":[{"id":"%s","name":"George","response":true,"has_responded":false}]}`, individual.ID))

		params := test.individualContract.GetInGroupFunc.History()[0]
		is.Equal(params.Arg1, individual.GroupID)
	})

	t.Run("can handle error", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		test.individualContract.GetInGroupFunc.PushReturn(nil, errors.New("query failed"))

		r := test.DoRequest(t, "GET", fmt.Sprintf("/api/v1/groups/%s/individuals", ced.NewID()), nil, http.StatusInternalServerError)
		is.Equal(r, `{"error":"Internal error","code":"internal"}`)
	})
}

func TestUpdateIndividual(t *testing.T) {
	individualID := ced.NewID()

	t.Run("param must be valid", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		r := test.DoRequest(t, "PUT", "/api/v1/individuals/--", nil, http.StatusUnprocessableEntity)
		is.Equal(r, `{"error":"invalid id","code":"unprocessable"}`)
	})

	t.Run("can update", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		test.individualContract.SetResponseFunc.PushReturn(nil)

		r := test.DoRequest(t, "PUT", fmt.Sprintf("/api/v1/individuals/%s", individualID), `{"response":true}`, http.StatusOK)
		is.Equal(r, "")

		params := test.individualContract.SetResponseFunc.History()[0]
		is.Equal(params.Arg1, individualID)
		is.Equal(params.Arg2, true)
	})

	t.Run("can update with blank request body", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		test.individualContract.SetResponseFunc.PushReturn(nil)

		r := test.DoRequest(t, "PUT", fmt.Sprintf("/api/v1/individuals/%s", individualID), nil, http.StatusOK)
		is.Equal(r, "")

		params := test.individualContract.SetResponseFunc.History()[0]
		is.Equal(params.Arg1, individualID)
		is.Equal(params.Arg2, false)
	})

	t.Run("handles parse error", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		r := test.DoRequest(t, "PUT", fmt.Sprintf("/api/v1/individuals/%s", ced.NewID()), `{"response":"shh"}`, http.StatusUnprocessableEntity)
		is.True(strings.Contains(r, `"code":"unprocessable"`))
	})

	t.Run("can handle error", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		test.individualContract.SetResponseFunc.PushReturn(errors.New("query failed"))

		r := test.DoRequest(t, "PUT", fmt.Sprintf("/api/v1/individuals/%s", ced.NewID()), nil, http.StatusInternalServerError)
		is.Equal(r, `{"error":"Internal error","code":"internal"}`)
	})
}
