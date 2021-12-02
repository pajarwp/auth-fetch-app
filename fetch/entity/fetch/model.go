package fetch

import "time"

type UserClaims struct {
	Phone     string `json:"phone"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type Resource struct {
	UUID          string `json:"uuid"`
	Komoditas     string `json:"komoditas"`
	AreaProvinsi  string `json:"area_provinsi"`
	AreaKota      string `json:"area_kota"`
	Size          string `json:"size"`
	Price         string `json:"price"`
	PriceUSD      string `json:"price_usd"`
	TanggalParsed string `json:"tgl_parsed"`
	Timestamp     string `json:"timestamp"`
}

type TempAggregate struct {
	TotalPrice float64
	ListPrice  []float64
	TotalSize  float64
	ListSize   []float64
	WeekStart  time.Time
	WeekEnd    time.Time
}
type Aggregate struct {
	Period      string    `json:"period"`
	MinSize     float64   `json:"min_size"`
	MaxSize     float64   `json:"max_size"`
	AvgSize     float64   `json:"avg_size"`
	MedianSize  float64   `json:"median_size"`
	ListSize    []float64 `json:"list_size"`
	MinPrice    float64   `json:"min_price"`
	MaxPrice    float64   `json:"max_price"`
	AvgPrice    float64   `json:"avg_price"`
	MedianPrice float64   `json:"median_price"`
	ListPrice   []float64 `json:"list_price"`
}

type Response struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}
