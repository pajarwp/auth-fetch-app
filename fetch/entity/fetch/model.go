package fetch

type UserClaims struct {
	Phone     string `json:"phone"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
