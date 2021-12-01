package main

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pajarwp/auth-fetch-app/entity/fetch/delivery"
	"github.com/pajarwp/auth-fetch-app/entity/fetch/repository"
	"github.com/pajarwp/auth-fetch-app/entity/fetch/usecase"
	"github.com/patrickmn/go-cache"
)

func main() {
	e := echo.New()
	c := cache.New(1*time.Hour, 1*time.Hour)
	repo := repository.NewAuthAppFetchRepository()
	usecase := usecase.NewFetchUsecase(c, repo)
	delivery.NewFetchHttpDelivery(*e, usecase)

	e.Start(":9000")
}
