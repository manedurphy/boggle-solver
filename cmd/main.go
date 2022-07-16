package main

import (
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"

	service "github.com/manedurphy/boggle-solver/internal"
)

func main() {
	var (
		app    *echo.Echo
		svc    *service.Service
		logger hclog.Logger
		err    error
	)

	logger = hclog.New(&hclog.LoggerOptions{
		Name:       "boggle service",
		Level:      hclog.Debug,
		Color:      hclog.ColorOff,
		TimeFormat: time.RFC3339Nano,
	})

	app = echo.New()
	svc = service.New(service.Config{
		Logger: logger,
	})

	app.POST("/solve", svc.Solve)
	if err = app.Start(":8080"); err != nil {
		logger.With("err", err).Error("could not start server")
		os.Exit(1)
	}
}
