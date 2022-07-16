package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/labstack/echo/v4"

	service "github.com/manedurphy/boggle-solver/internal"
)

var (
	cfg Config
)

type Config struct {
	Host         string `hcl:"host"`
	Port         int    `hcl:"port"`
	WordsZipFile string `hcl:"words_zip_file"`
	LogLevel     string `hcl:"log_level"`
}

func init() {
	if err := hclsimple.DecodeFile("configs/config.hcl", nil, &cfg); err != nil {
		panic(err)
	}
}

func main() {
	var (
		app    *echo.Echo
		svc    *service.Service
		logger hclog.Logger
		err    error
	)

	logger = hclog.New(&hclog.LoggerOptions{
		Name:       "boggle service",
		Level:      hclog.LevelFromString(cfg.LogLevel),
		Color:      hclog.ColorOff,
		TimeFormat: time.RFC3339Nano,
	})

	app = echo.New()
	svc = service.New(service.Config{
		Logger:       logger,
		WordsZipFile: cfg.WordsZipFile,
	})

	app.POST("/solve", svc.Solve)
	if err = app.Start(net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))); err != nil {
		logger.With("err", err).Error("could not start server")
		os.Exit(1)
	}
}
