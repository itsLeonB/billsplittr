package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

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

type JWTClaims struct {
	jwt.RegisteredClaims
	Data AuthenticatedData `json:"data"`
}

type AuthenticatedData struct {
	UserID uuid.UUID `json:"userId"`
}
