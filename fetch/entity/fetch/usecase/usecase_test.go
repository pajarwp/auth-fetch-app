package usecase_test

import (
	"testing"
	"time"

	"github.com/pajarwp/auth-fetch-app/entity/fetch/mock"
	"github.com/pajarwp/auth-fetch-app/entity/fetch/usecase"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

var usecaseInstance usecase.FetchUsecase

func init() {
	repoMock := mock.NewMockFetchRepository()
	usecaseInstance = usecase.NewFetchUsecase(cache.New(1*time.Hour, 1*time.Hour), repoMock)
}
func TestGetClaims(t *testing.T) {
	claims, err := usecaseInstance.GetClaims("token")
	assert.NoError(t, err)
	assert.Equal(t, "Nama Test", claims.Name)
}

func TestFetchresource(t *testing.T) {
	resource, err := usecaseInstance.FetchResource("token")
	assert.NoError(t, err)
	assert.Equal(t, "ACEH", resource[0].AreaProvinsi)
}

func TestAggregateResource(t *testing.T) {
	aggregate, err := usecaseInstance.AggregateResource("token")
	assert.NoError(t, err)
	assert.Equal(t, float64(10), aggregate["aceh"][0].MaxPrice)
}
