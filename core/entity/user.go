package entity

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required,min=6,max=15"`
	Password string `json:"password" validate:"required,numeric,len=6"`
	Role     string `json:"role" validate:"required,oneof=admin cashier"`
}

type RegisterResponse struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
