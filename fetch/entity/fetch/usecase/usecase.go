package usecase

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pajarwp/auth-fetch-app/entity/fetch"
	repo "github.com/pajarwp/auth-fetch-app/entity/fetch/repository"
	"github.com/patrickmn/go-cache"
)

type FetchUsecase interface {
	GetClaims(token string) (fetch.UserClaims, error)
	FetchResource(token string) ([]fetch.Resource, error)
}

type fetchUsecase struct {
	repo.FetchRepository
	*cache.Cache
}

func NewFetchUsecase(c *cache.Cache, repo repo.FetchRepository) FetchUsecase {
	return fetchUsecase{
		repo,
		c,
	}
}

func (f fetchUsecase) GetClaims(token string) (fetch.UserClaims, error) {
	return f.FetchRepository.GetClaims(token)
}

func (f fetchUsecase) FetchResource(token string) ([]fetch.Resource, error) {
	var listresource []fetch.Resource
	resp, err := f.FetchRepository.FetchResource(token)
	if err != nil {
		return listresource, err
	}

	converter, err := f.getCurrencyConverter()
	if err != nil {
		return listresource, err
	}

	for _, val := range resp {
		var resource fetch.Resource
		resource.UUID = val["uuid"]
		resource.Komoditas = val["komoditas"]
		resource.AreaProvinsi = val["area_provinsi"]
		resource.AreaKota = val["area_kota"]
		resource.Size = val["size"]
		resource.Price = val["price"]
		resource.TanggalParsed = val["tgl_parsed"]
		resource.Timestamp = val["timestamp"]
		if resource.Price != "" {
			price, err := strconv.ParseFloat(resource.Price, 64)
			if err != nil {
				return listresource, err
			}
			priceUSD := price * converter
			resource.PriceUSD = fmt.Sprintf("%0f", priceUSD)
		}
		listresource = append(listresource, resource)
	}
	return listresource, err
}

func (f fetchUsecase) getCurrencyConverter() (float64, error) {
	converterCache, found := f.Cache.Get("IDR_USD")
	if found {
		return converterCache.(float64), nil
	}
	converter, err := f.FetchRepository.GetCurrencyConverter()
	if err != nil {
		return 0, err
	}
	f.Cache.Set("IDR_USD", converter["IDR_USD"], 1*time.Hour)
	return converter["IDR_USD"], nil
}
