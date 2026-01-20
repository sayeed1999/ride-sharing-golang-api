package dto

type CustomerSignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func NewCustomerSignupRequest(email, name, password string) *CustomerSignupRequest {
	return &CustomerSignupRequest{
		Email:    email,
		Name:     name,
		Password: password,
	}
}
