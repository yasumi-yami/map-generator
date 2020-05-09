package main

import (
	"map-generator/infra/handler"
)

func main() {
	e := handler.Router()
	e.Logger.Fatal(e.Start(":8080"))
}
