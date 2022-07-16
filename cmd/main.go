package main

import (
	"github.com/labstack/echo/v4"

	service "github.com/manedurphy/boggle-solver/internal"
)

func main() {
	var (
		app *echo.Echo
		svc *service.Service
	)

	app = echo.New()
	svc = service.New()

	app.GET("/", svc.Healthz)
	app.POST("/solve", svc.Solve)
	panic(app.Start(":8080"))
}
