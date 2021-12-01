package usecase

import (
	"github.com/pajarwp/auth-fetch-app/entity/fetch"
	repo "github.com/pajarwp/auth-fetch-app/entity/fetch/repository"
)

type FetchUsecase interface {
	GetClaims(token string) (fetch.UserClaims, error)
}

type fetchUsecase struct {
	repo.FetchRepository
}

func NewFetchUsecase(repo repo.FetchRepository) FetchUsecase {
	return fetchUsecase{
		repo,
	}
}

func (f fetchUsecase) GetClaims(token string) (fetch.UserClaims, error) {
	return f.FetchRepository.GetClaims(token)
}
