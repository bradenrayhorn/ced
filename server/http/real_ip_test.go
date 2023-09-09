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
		xRealIP            string
		cfConnectingIP     string
		configTrustedIP    string
		expectedRemoteAddr string
	}{
		{"returns remote addr with no config",
			"192.0.0.1", "192.0.0.2", "", "::1"},
		{"returns remote addr with config and non-trusted cf ip",
			"192.0.0.1", "192.0.0.2", "192.0.0.3", "::1"},
		{"returns x-real-ip with config and trusted cf ip",
			"192.0.0.1", "192.0.0.2", "192.0.0.2", "192.0.0.1"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)

			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Add("X-Real-IP", test.xRealIP)
			req.Header.Add("CF-Connecting-IP", test.cfConnectingIP)
			req.RemoteAddr = "[::1]:65713"
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Use(RealIP(ced.Config{CloudflareTrustedIP: test.configTrustedIP}))

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
