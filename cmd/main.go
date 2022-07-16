package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/labstack/echo/v4"
	"github.com/manedurphy/boggle-solver/pkg/boggle"

	service "github.com/manedurphy/boggle-solver/internal"
)

var (
	cfg        Config
	configFile string
)

type Config struct {
	Host         string `hcl:"host"`
	Port         int    `hcl:"port"`
	WordsZipFile string `hcl:"words_zip_file"`
	LogLevel     string `hcl:"log_level"`
}

func init() {
	flag.StringVar(&configFile, "config-file", "configs/config.hcl", "path to the configuration file")
	flag.Parse()

	err := hclsimple.DecodeFile("configs/config.hcl", nil, &cfg)
	if err != nil {
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

	if err = boggle.InitDictionary(cfg.WordsZipFile); err != nil {
		logger.With("err", err).Error("failed to initialize dictionary")
		os.Exit(1)
	}

	app = echo.New()
	svc = service.New(service.Config{
		Logger:       logger,
		WordsZipFile: cfg.WordsZipFile,
	})

	app.GET("/", svc.Healthz)
	app.POST("/solve", svc.Solve)

	if err = app.Start(net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))); err != nil {
		logger.With("err", err).Error("could not start server")
		os.Exit(1)
	}
}
