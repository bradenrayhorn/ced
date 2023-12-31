package http

import (
	"net/http"
	"strings"

	"github.com/bradenrayhorn/ced/server/ced"
	"github.com/bradenrayhorn/ced/server/http/responses"
	"github.com/go-chi/chi/v5"
)

func (s *Server) handleGroupSearch() HandlerFunc {
	type response []responses.Group

	return func(r *http.Request) (HttpResponse, error) {
		search := strings.TrimSpace(r.URL.Query().Get("search"))

		if search == "" {
			return HttpResponse{}, ced.NewError(ced.EUNPROCESSABLE, "search is required")
		}

		groups, err := s.groupContract.Search(r.Context(), makeReqCtx(r), search)
		if err != nil {
			return HttpResponse{}, err
		}

		return HttpResponse{status: http.StatusOK, body: struct {
			Data response `json:"data"`
		}{
			Data: responses.MapSlice(groups, responses.FromGroup),
		}}, nil
	}
}

func (s *Server) handleGroupGet() HandlerFunc {
	return func(r *http.Request) (HttpResponse, error) {
		id, err := ced.IDFromString(chi.URLParam(r, "groupID"))

		if err != nil {
			return HttpResponse{}, ced.NewError(ced.EUNPROCESSABLE, "invalid id")
		}

		group, err := s.groupContract.Get(r.Context(), makeReqCtx(r), id)
		if err != nil {
			return HttpResponse{}, err
		}

		return HttpResponse{status: http.StatusOK, body: struct {
			Data responses.Group `json:"data"`
		}{
			Data: responses.FromGroup(group),
		}}, nil
	}
}

func (s *Server) handleGroupUpdate() HandlerFunc {
	type request struct {
		Attendees uint8 `json:"attendees,string"`
	}

	return func(r *http.Request) (HttpResponse, error) {
		id, err := ced.IDFromString(chi.URLParam(r, "groupID"))

		if err != nil {
			return HttpResponse{}, ced.NewError(ced.EUNPROCESSABLE, "invalid id")
		}

		var req request
		if err := decodeRequest(r, &req); err != nil {
			return HttpResponse{}, err
		}

		if err := s.groupContract.Respond(r.Context(), makeReqCtx(r), id, req.Attendees); err != nil {
			return HttpResponse{}, err
		}

		return HttpResponse{status: http.StatusOK, body: nil}, nil
	}
}
