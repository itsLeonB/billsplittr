package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/billsplittr/internal/test/mocks_test"
	"github.com/itsLeonB/cocoon-protos/gen/go/auth/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks_test.NewMockAuthServiceClient(ctrl)
	authService := service.NewAuthService(mockClient)

	request := dto.RegisterRequest{
		Email:                "test@example.com",
		Password:             "password123",
		PasswordConfirmation: "password123",
	}

	mockClient.EXPECT().
		Register(gomock.Any(), &auth.RegisterRequest{
			Email:                request.Email,
			Password:             request.Password,
			PasswordConfirmation: request.PasswordConfirmation,
		}).
		Return(&auth.RegisterResponse{}, nil)

	err := authService.Register(context.Background(), request)
	assert.NoError(t, err)
}

func TestAuthService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks_test.NewMockAuthServiceClient(ctrl)
	authService := service.NewAuthService(mockClient)

	request := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedResponse := &auth.LoginResponse{
		Type:  "Bearer",
		Token: "test-token",
	}

	mockClient.EXPECT().
		Login(gomock.Any(), &auth.LoginRequest{
			Email:    request.Email,
			Password: request.Password,
		}).
		Return(expectedResponse, nil)

	response, err := authService.Login(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, "Bearer", response.Type)
	assert.Equal(t, "test-token", response.Token)
}

func TestAuthService_VerifyToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks_test.NewMockAuthServiceClient(ctrl)
	authService := service.NewAuthService(mockClient)

	token := "valid-token"
	profileID := "123e4567-e89b-12d3-a456-426614174000"

	expectedResponse := &auth.VerifyTokenResponse{
		ProfileId: profileID,
	}

	mockClient.EXPECT().
		VerifyToken(gomock.Any(), &auth.VerifyTokenRequest{Token: token}).
		Return(expectedResponse, nil)

	valid, data, err := authService.VerifyToken(context.Background(), token)

	assert.NoError(t, err)
	assert.True(t, valid)
	assert.Equal(t, profileID, data[appconstant.ContextProfileID])
}

func TestAuthService_VerifyToken_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks_test.NewMockAuthServiceClient(ctrl)
	authService := service.NewAuthService(mockClient)

	token := "invalid-token"

	mockClient.EXPECT().
		VerifyToken(gomock.Any(), &auth.VerifyTokenRequest{Token: token}).
		Return(nil, status.Error(codes.Unauthenticated, "unauthorized"))

	valid, data, err := authService.VerifyToken(context.Background(), token)

	assert.Error(t, err)
	assert.False(t, valid)
	assert.Nil(t, data)
}
