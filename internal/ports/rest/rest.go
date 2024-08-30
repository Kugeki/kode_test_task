package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type Server struct {
	s   *http.Server
	log *slog.Logger
}

//	@title			Note Service
//	@version		0.1.0
//	@description	Create notes. Get your notes.

// @BasePath	/

func NewServer(router chi.Router, log *slog.Logger, options ...Opt) (*Server, error) {
	s := &http.Server{Handler: router}

	for _, op := range options {
		err := op(s)
		if err != nil {
			return nil, err
		}
	}

	log = log.With(slog.String("context", "rest server"))
	return &Server{s: s, log: log}, nil
}

func (s *Server) Run() error {
	s.log.Info("listening", slog.String("addr", s.s.Addr))
	return s.s.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("shutdown is started")
	return s.s.Shutdown(ctx)
}
