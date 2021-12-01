package main

import (
	"github.com/labstack/echo/v4"
	"github.com/pajarwp/auth-fetch-app/entity/fetch/delivery"
	"github.com/pajarwp/auth-fetch-app/entity/fetch/repository"
	"github.com/pajarwp/auth-fetch-app/entity/fetch/usecase"
)

func main() {
	e := echo.New()
	repo := repository.NewAuthAppFetchRepository()
	usecase := usecase.NewFetchUsecase(repo)
	delivery.NewFetchHttpDelivery(*e, usecase)

	e.Start(":9000")
}
