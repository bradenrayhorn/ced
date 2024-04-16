package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/matryer/is"
)

func TestRealIP(t *testing.T) {
	var tests = []struct {
		name               string
		connectingIPHeader string
		expectedRemoteAddr string
	}{
		{"returns remote addr and no value in header",
			"", "::1"},
		{"returns real ip and value in header",
			"192.0.0.1", "192.0.0.1"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)

			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Add("ced-connecting-ip", test.connectingIPHeader)
			req.RemoteAddr = "[::1]:65713"
			w := httptest.NewRecorder()

			r := chi.NewRouter()

			r.Use(RealIP())

			realIP := ""
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				realIP = r.RemoteAddr
			})
			r.ServeHTTP(w, req)
			is.Equal(w.Result().StatusCode, 200)

			is.Equal(test.expectedRemoteAddr, realIP)
		})
	}
}
