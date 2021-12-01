package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(ctx echo.Context) error {
		data := "Hello World"
		return ctx.String(http.StatusOK, data)
	})

	e.Start(":9000")
}
