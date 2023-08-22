package http

import (
	"encoding/json"
	"net/http"
)

func writeJsonResponse(w http.ResponseWriter, v any, statusCode int) {
	if v != nil {
		w.Header().Add("content-type", "application/json")
	}
	w.WriteHeader(statusCode)

	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			writeError(w, err)
		}
	}
}
