package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradenrayhorn/ced/server/ced"
	"github.com/bradenrayhorn/ced/server/internal/testutils"
	"github.com/matryer/is"
)

func TestDecodeRequest(t *testing.T) {
	request := func(body string) *http.Request {
		return httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
	}

	t.Run("does not fail with empty body", func(t *testing.T) {
		is := is.New(t)
		req := request(``)
		var r map[string]interface{}
		err := decodeRequest(req, &r)

		is.NoErr(err)
	})

	t.Run("returns error for bad syntax", func(t *testing.T) {
		req := request(`{"name":}`)
		var r map[string]interface{}
		err := decodeRequest(req, &r)

		testutils.IsCode(t, err, ced.EUNPROCESSABLE)
	})

	t.Run("returns error with missing quotes syntax", func(t *testing.T) {
		req := request(`{"name:}`)
		var r map[string]interface{}
		err := decodeRequest(req, &r)

		testutils.IsCode(t, err, ced.EUNPROCESSABLE)
	})

	t.Run("returns error for bad type", func(t *testing.T) {
		req := request(`{"name":"test"}`)
		var r struct {
			Name bool `json:"name"`
		}
		err := decodeRequest(req, &r)

		testutils.IsCodeAndError(t, err, ced.EUNPROCESSABLE, "Invalid data `string` for `name` field.")
	})

	t.Run("returns error for bad type with ,string flag", func(t *testing.T) {
		req := request(`{"val":false}`)
		var r struct {
			Val bool `json:"val,string"`
		}
		err := decodeRequest(req, &r)

		testutils.IsCodeAndError(t, err, ced.EUNPROCESSABLE, "Invalid data provided. Please double check data types.")
	})
}
