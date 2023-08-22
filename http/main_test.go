package http_test

import (
	"testing"

	"github.com/bradenrayhorn/ced/http"
	"github.com/bradenrayhorn/ced/internal/mocks"
	"github.com/bradenrayhorn/ced/internal/testutils"
)

type httpTest struct {
	httpServer *http.Server

	individualContract *mocks.MockIndividualContract
}

func newHttpTest() *httpTest {
	test := &httpTest{
		individualContract: mocks.NewMockIndividualContract(),
	}

	test.httpServer = http.NewServer(
		test.individualContract,
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
	return testutils.DoRequest(tb, method, "http://"+t.httpServer.GetBoundAddr()+path, body, expectedStatus)
}
