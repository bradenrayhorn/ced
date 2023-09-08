package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/bradenrayhorn/ced/server/ced"
)

var codeToHTTPStatus = map[string]int{
	ced.EFORBIDDEN:     http.StatusForbidden,
	ced.EINTERNAL:      http.StatusInternalServerError,
	ced.EINVALID:       http.StatusUnprocessableEntity,
	ced.ENOTFOUND:      http.StatusNotFound,
	ced.EUNAUTHORIZED:  http.StatusUnauthorized,
	ced.EUNPROCESSABLE: http.StatusBadRequest,
}

func writeError(w http.ResponseWriter, err error) {
	type errorResponse struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}

	var code = ced.EINTERNAL
	var msg = ced.ErrorInternal.Error()
	var cedError ced.Error
	if errors.As(err, &cedError) {
		code, msg = cedError.CedError()
	}

	if code == ced.EINTERNAL {
		slog.Error(err.Error())
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(codeToHTTPStatus[code])
	_ = json.NewEncoder(w).Encode(&errorResponse{Code: code, Error: msg})
}
