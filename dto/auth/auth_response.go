package authdto

type LoginResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Number  string `json:"number"`
	Address string `json:"address"`
	Role    string `json:"role"`
	Token   string `json:"token"`
}
