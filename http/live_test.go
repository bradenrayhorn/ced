package http_test

import (
	"net/http"
	"testing"

	"github.com/matryer/is"
)

func TestLive(t *testing.T) {
	test := newHttpTest()
	defer test.Stop(t)
	is := is.New(t)

	r := test.DoRequest(t, "GET", "/api/v1/live", nil, http.StatusOK)
	is.Equal(r, "ok")
}
