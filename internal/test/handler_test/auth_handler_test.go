package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/billsplittr/internal/delivery/http/handler"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_HandleRegister_InvalidJSON(t *testing.T) {
	authService := &mockAuthService{}
	authHandler := handler.NewAuthHandler(authService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/register", authHandler.HandleRegister())

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusCreated, w.Code)
}

type mockAuthService struct{}

func (m *mockAuthService) Register(ctx context.Context, request dto.RegisterRequest) error {
	return nil
}

func (m *mockAuthService) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error) {
	return dto.LoginResponse{Type: "Bearer", Token: "test-token"}, nil
}

func (m *mockAuthService) VerifyToken(ctx context.Context, token string) (bool, map[string]any, error) {
	return true, map[string]any{"profile_id": "test-id"}, nil
}
