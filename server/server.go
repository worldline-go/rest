package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/worldline-go/logz/logecho"
	"github.com/worldline-go/tell/metric/metricecho"
	"github.com/ziflex/lecho/v3"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var ShutdownTimeout = 10 * time.Second

type Server struct {
	mux    *http.ServeMux
	server *http.Server

	m sync.Mutex
}

type ServerFunc func(ctx context.Context, mux *http.ServeMux, e *echo.Echo) error

func (s *Server) Start(addr string) error {
	s.server = &http.Server{ //nolint:gosec // skip check in service
		Addr:    addr,
		Handler: h2c.NewHandler(s.mux, &http2.Server{}),
	}

	log.Info().Msgf("starting server on port %s", s.server.Addr)

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

	log.Warn().Msg("stopping service")
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

func New(ctx context.Context, name string, fn func(ctx context.Context, mux *http.ServeMux, e *echo.Echo) error) (*Server, error) {
	mux := http.NewServeMux()

	// echo server
	e := echo.New()
	e.HideBanner = true

	e.Validator = NewValidator()

	e.Logger = lecho.From(log.Logger)

	e.Use(
		middleware.Recover(),
		middleware.CORS(),
		middleware.RequestID(),
		middleware.RequestLoggerWithConfig(logecho.RequestLoggerConfig()),
		logecho.ZerologLogger(),
		metricecho.HTTPMetrics(),
		otelecho.Middleware(name),
		MiddlewareUserInfo,
	)

	e.HTTPErrorHandler = HTTPErrorHandler

	if err := fn(ctx, mux, e); err != nil {
		return nil, fmt.Errorf("failed to register routes: %w", err)
	}

	return &Server{
		mux: mux,
	}, nil
}
