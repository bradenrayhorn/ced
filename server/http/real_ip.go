package http

import (
	"net"
	"net/http"

	"github.com/bradenrayhorn/ced/server/ced"
)

func RealIP(config ced.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if ip := getRealIP(config, req); ip != "" {
				req.RemoteAddr = ip
			}
			next.ServeHTTP(w, req)
		})
	}
}

func getRealIP(config ced.Config, req *http.Request) string {
	var realIP = ""
	// default - parse remote addr
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		realIP = ip
	}

	// cloudflare
	if config.CloudflareTrustedIP != "" {
		if req.Header.Get("CF-Connecting-IP") == config.CloudflareTrustedIP {
			realIP = req.Header.Get("x-ced-connecting-ip")
		}
	}

	return realIP
}
