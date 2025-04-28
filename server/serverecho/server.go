package serverecho

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/worldline-go/logz/logecho"
	"github.com/worldline-go/tell/metric/metricecho"
	"github.com/ziflex/lecho/v3"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func New(ctx context.Context, name string) (*echo.Echo, error) {
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

	return e, nil
}
