package authdto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password" validate:"required"`
	Number   string `json:"number"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"`
}
