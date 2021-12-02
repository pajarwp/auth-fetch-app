package mock

import (
	"github.com/pajarwp/auth-fetch-app/entity/fetch"
)

type MockFetchRepository interface {
	GetClaims(token string) (fetch.UserClaims, error)
	FetchResource(token string) ([]map[string]string, error)
	GetCurrencyConverter() (map[string]float64, error)
}

type mockFetchRepository struct {
}

func NewMockFetchRepository() MockFetchRepository {
	return mockFetchRepository{}
}

func (a mockFetchRepository) GetClaims(token string) (fetch.UserClaims, error) {
	claims := fetch.UserClaims{}
	claims.Name = "Nama Test"
	claims.Phone = "081234567890"
	claims.Role = "admin"
	claims.CreatedAt = "2021-12-02 21:36:29"
	return claims, nil
}

func (a mockFetchRepository) FetchResource(token string) ([]map[string]string, error) {
	respMap := []map[string]string{
		{
			"uuid":          "8bf8d336-47c0-447d-a7ba-5fb85537cef8",
			"komoditas":     "Naga",
			"area_provinsi": "ACEH",
			"area_kota":     "ACEH KOTA",
			"size":          "30",
			"price":         "10",
			"tgl_parsed":    "2021-10-28T15:44:15.652Z",
			"timestamp":     "1635435855652",
		},
		{
			"uuid":          "f0db7091-b4c4-4b51-80f4-541a526f08aa",
			"komoditas":     "Salmonela",
			"area_provinsi": "BANTEN",
			"area_kota":     "PANDEGLANG",
			"size":          "50",
			"price":         "10",
			"tgl_parsed":    "2021-10-28T15:46:34.871Z",
			"timestamp":     "1635435994871",
		},
	}
	return respMap, nil
}

func (a mockFetchRepository) GetCurrencyConverter() (map[string]float64, error) {
	respMap := map[string]float64{
		"IDR_USD": 0.1,
	}

	return respMap, nil
}
