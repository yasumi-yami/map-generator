package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const indent = "  "

// Router ルーティング
func Router() *echo.Echo {
	e := echo.New()

	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/", func(c echo.Context) error {
		return c.JSONPretty(
			http.StatusOK,
			map[string]interface{}{"text": "hello yami world"},
			indent,
		)
	})

	mapGeneratorAPI := NewMapGeneratorAPI()
	mapGenerator := e.Group("/map")
	mapGenerator.GET("", mapGeneratorAPI.Generate)
	mapGenerator.GET("/:id", mapGeneratorAPI.Get)
	return e
}
