package http

import (
	"net/http"
	"strings"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/http/responses"
	"github.com/go-chi/chi/v5"
)

func (s *Server) handleIndividualSearch() HandlerFunc {
	type response map[string][]responses.Individual

	return func(r *http.Request) (HttpResponse, error) {
		search := strings.TrimSpace(r.URL.Query().Get("search"))

		if search == "" {
			return HttpResponse{}, ced.NewError(ced.EUNPROCESSABLE, "search is required")
		}

		groups, err := s.individualContract.SearchByName(r.Context(), search)
		if err != nil {
			return HttpResponse{}, err
		}

		res := make(response)
		for k, v := range groups {
			res[k.String()] = responses.MapSlice(v, responses.FromIndividual)
		}

		return HttpResponse{status: http.StatusOK, body: struct {
			Data response `json:"data"`
		}{
			Data: res,
		}}, nil
	}
}

func (s *Server) handleGetIndividualsInGroup() HandlerFunc {
	type response []responses.Individual

	return func(r *http.Request) (HttpResponse, error) {
		groupID, err := ced.IDFromString(chi.URLParam(r, "groupID"))

		if err != nil {
			return HttpResponse{}, ced.NewError(ced.EUNPROCESSABLE, "invalid group id")
		}

		individuals, err := s.individualContract.GetInGroup(r.Context(), groupID)
		if err != nil {
			return HttpResponse{}, err
		}

		return HttpResponse{status: http.StatusOK, body: struct {
			Data response `json:"data"`
		}{
			Data: responses.MapSlice(individuals, responses.FromIndividual),
		}}, nil
	}
}

func (s *Server) handleIndividualUpdate() HandlerFunc {
	type request struct {
		Response bool `json:"response"`
	}

	return func(r *http.Request) (HttpResponse, error) {
		id, err := ced.IDFromString(chi.URLParam(r, "individualID"))

		if err != nil {
			return HttpResponse{}, ced.NewError(ced.EUNPROCESSABLE, "invalid id")
		}

		var req request
		if err := decodeRequest(r, &req); err != nil {
			return HttpResponse{}, err
		}

		if err := s.individualContract.SetResponse(r.Context(), id, req.Response); err != nil {
			return HttpResponse{}, err
		}

		return HttpResponse{status: http.StatusOK, body: nil}, nil
	}
}
