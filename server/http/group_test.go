package http_test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/bradenrayhorn/ced/server/ced"
	"github.com/matryer/is"
)

func TestGroupSearch(t *testing.T) {
	t.Run("param is required", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		r := test.DoRequest(t, "GET", "/api/v1/groups/search", nil, http.StatusUnprocessableEntity)
		is.Equal(r, `{"error":"search is required","code":"unprocessable"}`)
	})

	t.Run("can search", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		group := ced.Group{
			ID:           ced.NewID(),
			Name:         "George",
			Attendees:    1,
			MaxAttendees: 1,
			HasResponded: true,
		}
		test.groupContract.SearchFunc.PushReturn([]ced.Group{group}, nil)

		r := test.DoRequest(t, "GET", "/api/v1/groups/search?search=geo", nil, http.StatusOK)
		is.Equal(r, fmt.Sprintf(`{"data":[{"id":"%s","name":"George","attendees":1,"max_attendees":1,"has_responded":true}]}`, group.ID))

		params := test.groupContract.SearchFunc.History()[0]
		is.Equal(params.Arg1, "geo")
	})

	t.Run("can handle search error", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		test.groupContract.SearchFunc.PushReturn(nil, errors.New("query failed"))

		r := test.DoRequest(t, "GET", "/api/v1/groups/search?search=geo", nil, http.StatusInternalServerError)
		is.Equal(r, `{"error":"Internal error","code":"internal"}`)
	})
}

func TestGetGroup(t *testing.T) {
	t.Run("param must be valid", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		r := test.DoRequest(t, "GET", "/api/v1/groups/--", nil, http.StatusUnprocessableEntity)
		is.Equal(r, `{"error":"invalid id","code":"unprocessable"}`)
	})

	t.Run("can get in group", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		group := ced.Group{
			ID:           ced.NewID(),
			Name:         "George",
			Attendees:    1,
			MaxAttendees: 1,
			HasResponded: true,
		}
		test.groupContract.GetFunc.PushReturn(group, nil)

		r := test.DoRequest(t, "GET", fmt.Sprintf("/api/v1/groups/%s", group.ID), nil, http.StatusOK)
		is.Equal(r, fmt.Sprintf(`{"data":{"id":"%s","name":"George","attendees":1,"max_attendees":1,"has_responded":true}}`, group.ID))

		params := test.groupContract.GetFunc.History()[0]
		is.Equal(params.Arg1, group.ID)
	})

	t.Run("can handle error", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		test.groupContract.GetFunc.PushReturn(ced.Group{}, errors.New("query failed"))

		r := test.DoRequest(t, "GET", fmt.Sprintf("/api/v1/groups/%s", ced.NewID()), nil, http.StatusInternalServerError)
		is.Equal(r, `{"error":"Internal error","code":"internal"}`)
	})
}

func TestUpdateGroup(t *testing.T) {
	groupID := ced.NewID()

	t.Run("param must be valid", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		r := test.DoRequest(t, "PUT", "/api/v1/groups/--", nil, http.StatusUnprocessableEntity)
		is.Equal(r, `{"error":"invalid id","code":"unprocessable"}`)
	})

	t.Run("can update", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		test.groupContract.RespondFunc.PushReturn(nil)

		r := test.DoRequest(t, "PUT", fmt.Sprintf("/api/v1/groups/%s", groupID), `{"attendees":5}`, http.StatusOK)
		is.Equal(r, "")

		params := test.groupContract.RespondFunc.History()[0]
		is.Equal(params.Arg1, groupID)
		is.Equal(params.Arg2, uint8(5))
	})

	t.Run("can update with blank request body", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		test.groupContract.RespondFunc.PushReturn(nil)

		r := test.DoRequest(t, "PUT", fmt.Sprintf("/api/v1/groups/%s", groupID), nil, http.StatusOK)
		is.Equal(r, "")

		params := test.groupContract.RespondFunc.History()[0]
		is.Equal(params.Arg1, groupID)
		is.Equal(params.Arg2, uint8(0))
	})

	t.Run("handles parse error", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		r := test.DoRequest(t, "PUT", fmt.Sprintf("/api/v1/groups/%s", ced.NewID()), `{"attendees":"shh"}`, http.StatusUnprocessableEntity)
		is.True(strings.Contains(r, `"code":"unprocessable"`))
	})

	t.Run("can handle error", func(t *testing.T) {
		test := newHttpTest()
		defer test.Stop(t)
		is := is.New(t)

		test.groupContract.RespondFunc.PushReturn(errors.New("query failed"))

		r := test.DoRequest(t, "PUT", fmt.Sprintf("/api/v1/groups/%s", ced.NewID()), nil, http.StatusInternalServerError)
		is.Equal(r, `{"error":"Internal error","code":"internal"}`)
	})
}
