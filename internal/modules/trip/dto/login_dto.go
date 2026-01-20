package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func NewLoginRequest(email, password string) *LoginRequest {
	return &LoginRequest{
		Email:    email,
		Password: password,
	}
}
