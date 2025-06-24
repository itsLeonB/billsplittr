package dto

type RegisterRequest struct {
	Email                string `json:"email" binding:"required,email,min=3"`
	Password             string `json:"password" binding:"required,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,min=3"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}
