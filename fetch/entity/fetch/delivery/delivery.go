package delivery

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pajarwp/auth-fetch-app/entity/fetch"
	usecase "github.com/pajarwp/auth-fetch-app/entity/fetch/usecase"
)

type FetchHttpDelivery struct {
	usecase.FetchUsecase
}

func NewFetchHttpDelivery(e echo.Echo, usecase usecase.FetchUsecase) {
	handler := FetchHttpDelivery{
		usecase,
	}
	fetch := e.Group("/fetch")
	fetch.GET("/claims", handler.GetClaims)
	fetch.GET("/resource", handler.FetchResource)
	fetch.GET("/aggregate", handler.AggregateResource)
}

func (f FetchHttpDelivery) GetClaims(c echo.Context) error {
	resp := fetch.Response{}
	token := c.Request().Header.Get(echo.HeaderAuthorization)
	if token == "" {
		resp.Status = "Unauthorized"
		resp.Message = "Missing Authorization Token in Header"
		return c.JSON(http.StatusUnauthorized, resp)
	}
	claims, err := f.FetchUsecase.GetClaims(token)
	if err != nil {
		resp.Status = "failed"
		resp.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Status = "ok"
	resp.Message = "Success"
	resp.Data = claims
	return c.JSON(http.StatusOK, resp)
}

func (f FetchHttpDelivery) FetchResource(c echo.Context) error {
	resp := fetch.Response{}
	token := c.Request().Header.Get(echo.HeaderAuthorization)
	if token == "" {
		resp.Status = "Unauthorized"
		resp.Message = "Missing Authorization Token in Header"
		return c.JSON(http.StatusUnauthorized, resp)
	}
	resources, err := f.FetchUsecase.FetchResource(token)
	if err != nil {
		resp.Status = "failed"
		resp.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Status = "ok"
	resp.Message = "Success"
	resp.Data = resources
	return c.JSON(http.StatusOK, resp)
}

func (f FetchHttpDelivery) AggregateResource(c echo.Context) error {
	resp := fetch.Response{}
	token := c.Request().Header.Get(echo.HeaderAuthorization)
	if token == "" {
		resp.Status = "Unauthorized"
		resp.Message = "Missing Authorization Token in Header"
		return c.JSON(http.StatusUnauthorized, resp)
	}
	resources, err := f.FetchUsecase.AggregateResource(token)
	if err != nil {
		resp.Status = "failed"
		resp.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Status = "ok"
	resp.Message = "Success"
	resp.Data = resources
	return c.JSON(http.StatusOK, resp)
}
