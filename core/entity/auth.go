package entity

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=6,max=15"`
	Password string `json:"password" validate:"required,numeric,len=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
