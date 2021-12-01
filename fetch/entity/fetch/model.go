package fetch

type UserClaims struct {
	Phone     string `json:"phone"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type Response struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}
