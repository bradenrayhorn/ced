package http

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router    *chi.Mux
	sv        *http.Server
	boundAddr string

	individualContract ced.IndividualContract
}

func (s *Server) Open(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.boundAddr = ln.Addr().String()

	go func() {
		err := s.sv.Serve(ln)
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start http server", "error", err)
		}
	}()

	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	return s.sv.Shutdown(ctx)
}

func (s *Server) GetBoundAddr() string {
	return s.boundAddr
}

func NewServer(
	individualContract ced.IndividualContract,
) *Server {
	s := &Server{
		router: chi.NewRouter(),
		sv:     &http.Server{},

		individualContract: individualContract,
	}

	s.sv.Handler = s.router
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Get("/live", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("ok"))
		})

		r.Route("/individuals", func(r chi.Router) {
			r.Get("/search", toHttpHandlerFunc(s.handleIndividualSearch()))
			r.Put("/{individualID}", toHttpHandlerFunc(s.handleIndividualUpdate()))
		})

		r.Route("/groups", func(r chi.Router) {
			r.Get("/{groupID}/individuals", toHttpHandlerFunc(s.handleGetIndividualsInGroup()))
		})
	})

	return s
}
