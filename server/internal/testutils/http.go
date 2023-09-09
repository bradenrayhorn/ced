package testutils

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func DoRequest(t testing.TB, method string, url string, body any, expectedStatus int, prepare func(r *http.Request) *http.Request) string {
	var reqBody io.Reader
	switch body := body.(type) {
	case string:
		reqBody = bytes.NewReader([]byte(body))
	default:
		reqBody = nil
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req = prepare(req)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	return strings.TrimSpace(string(data))
}
