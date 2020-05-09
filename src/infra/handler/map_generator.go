package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"map-generator/application/service"
)

const INDENT = "  "

// MapGeneratorAPI 地図情報を作成する
type MapGeneratorAPI struct {
	generator *service.MapGenerator
}

// GetMap 地図情報の取得
func (handler *MapGeneratorAPI) GetMap(c echo.Context) error {
	res := handler.generator.Generate(10)
	return c.JSONPretty(http.StatusOK, res, INDENT)
}
