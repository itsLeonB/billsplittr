package dto

type RegisterRequest struct {
	Username             string `json:"username" binding:"required"`
	Password             string `json:"password" binding:"required,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}
