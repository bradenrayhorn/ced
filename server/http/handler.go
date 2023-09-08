package http

import "net/http"

type HttpResponse struct {
	status int
	body   any
}

type HandlerFunc func(r *http.Request) (HttpResponse, error)

func toHttpHandlerFunc(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := f(r)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJsonResponse(w, res.body, res.status)
	}
}
