package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"map-generator/application/service"
)

const INDENT = "  "

// MapGeneratorAPI 地図情報を作成する
type MapGeneratorAPI struct {
	generator *service.MapGenerator
}

func NewMapGeneratorAPI() *MapGeneratorAPI {
	return &MapGeneratorAPI{
		generator: &service.MapGenerator{},
	}
}

// Generate 地図情報の生成 TODO オンライン処理でなくする
func (api *MapGeneratorAPI) Generate(c echo.Context) error {
	res := api.generator.Generate(500)
	return c.JSONPretty(http.StatusOK, res, INDENT)
}

// Get 取得
func (api *MapGeneratorAPI) Get(c echo.Context) error {
	center, _ := strconv.Atoi(c.Param("id"))

	res := api.generator.Get(center)
	return c.JSONPretty(http.StatusOK, res, INDENT)
}
