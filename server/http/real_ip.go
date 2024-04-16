package http

import (
	"net"
	"net/http"
)

func RealIP() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if ip := getRealIP(req); ip != "" {
				req.RemoteAddr = ip
			}
			next.ServeHTTP(w, req)
		})
	}
}

func getRealIP(req *http.Request) string {
	var realIP = ""
	// default - parse remote addr
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		realIP = ip
	}

	// get ced-connecting-ip header (set by svelte)
	if len(req.Header.Get("ced-connecting-ip")) != 0 {
		realIP = req.Header.Get("ced-connecting-ip")
	}

	return realIP
}
