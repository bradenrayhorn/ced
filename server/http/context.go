package http

import (
	"net/http"

	"github.com/bradenrayhorn/ced/server/ced"
)

func makeReqCtx(r *http.Request) ced.ReqContext {
	return ced.ReqContext{ConnectingIP: r.RemoteAddr}
}
