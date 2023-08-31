package http

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	router    *chi.Mux
	sv        *http.Server
	boundAddr string

	groupContract ced.GroupContract
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
	config ced.Config,
	groupContract ced.GroupContract,
) *Server {
	s := &Server{
		router: chi.NewRouter(),
		sv:     &http.Server{},

		groupContract: groupContract,
	}

	s.sv.Handler = s.router
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins: strings.Split(config.AllowedOrigin, ","),
	}))

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Get("/live", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("ok"))
		})

		r.Route("/groups", func(r chi.Router) {
			r.Get("/search", toHttpHandlerFunc(s.handleGroupSearch()))

			r.Get("/{groupID}", toHttpHandlerFunc(s.handleGroupGet()))
			r.Put("/{groupID}", toHttpHandlerFunc(s.handleGroupUpdate()))
		})
	})

	return s
}
