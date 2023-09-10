package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradenrayhorn/ced/server/ced"
	"github.com/go-chi/chi/v5"
	"github.com/matryer/is"
)

func TestRealIP(t *testing.T) {
	var tests = []struct {
		name               string
		addConfig          bool
		xRealIP            string
		expectedRemoteAddr string
	}{
		{"returns remote addr with no config",
			false, "192.0.0.1", "::1"},
		{"returns remote addr with config and no value in header",
			true, "", "::1"},
		{"returns real ip with config and value in header",
			true, "192.0.0.1", "192.0.0.1"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)

			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Add("x-real-ip", test.xRealIP)
			req.RemoteAddr = "[::1]:65713"
			w := httptest.NewRecorder()

			r := chi.NewRouter()

			header := ""
			if test.addConfig {
				header = "x-real-ip"
			}
			r.Use(RealIP(ced.Config{TrustedClientIPHeader: header}))

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
