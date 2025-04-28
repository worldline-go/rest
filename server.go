package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var ShutdownTimeout = 10 * time.Second

type Server struct {
	Mux    *http.ServeMux
	server *http.Server

	m sync.Mutex
}

func New() *Server {
	return &Server{
		Mux: http.NewServeMux(),
	}
}

func (s *Server) Start(addr string) error {
	s.server = &http.Server{ //nolint:gosec // skip check in service
		Addr:    addr,
		Handler: h2c.NewHandler(s.Mux, &http2.Server{}),
	}

	log.Info().Msgf("starting server on %s", addr)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Stop() error {
	s.m.Lock()
	defer s.m.Unlock()

	if s.server == nil {
		return nil
	}

	log.Warn().Msg("stopping server")
	defer func() {
		s.server = nil
	}()

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown service: %w", err)
	}

	return nil
}
