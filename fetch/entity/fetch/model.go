package fetch

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

type Response struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}
