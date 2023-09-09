package http_test

import (
	"testing"

	goHttp "net/http"

	"github.com/bradenrayhorn/ced/server/ced"
	"github.com/bradenrayhorn/ced/server/http"
	"github.com/bradenrayhorn/ced/server/internal/mocks"
	"github.com/bradenrayhorn/ced/server/internal/testutils"
)

type httpTest struct {
	httpServer *http.Server

	groupContract *mocks.MockGroupContract
}

func newHttpTest() *httpTest {
	test := &httpTest{
		groupContract: mocks.NewMockGroupContract(),
	}

	test.httpServer = http.NewServer(
		ced.Config{
			CloudflareTrustedIP: "8.7.7.7",
		},
		test.groupContract,
	)

	if err := test.httpServer.Open(":0"); err != nil {
		panic(err)
	}

	return test
}

func (t *httpTest) Stop(tb testing.TB) {
	if err := t.httpServer.Close(); err != nil {
		tb.Fatal(err)
	}
}

func (t *httpTest) DoRequest(tb testing.TB, method string, path string, body any, expectedStatus int) string {
	return testutils.DoRequest(tb, method, "http://"+t.httpServer.GetBoundAddr()+path, body, expectedStatus, func(r *goHttp.Request) *goHttp.Request { return r })
}

func (t *httpTest) DoRequestWith(tb testing.TB, method string, path string, body any, expectedStatus int, prepare func(r *goHttp.Request) *goHttp.Request) string {
	return testutils.DoRequest(tb, method, "http://"+t.httpServer.GetBoundAddr()+path, body, expectedStatus, prepare)
}
