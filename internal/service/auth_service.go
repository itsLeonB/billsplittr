package service

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/cocoon-protos/gen/go/auth"
	"github.com/rotisserie/eris"
)

type authServiceGrpc struct {
	authClient auth.AuthServiceClient
}

func NewAuthService(
	authClient auth.AuthServiceClient,
) AuthService {
	return &authServiceGrpc{
		authClient,
	}
}

func (as *authServiceGrpc) Register(ctx context.Context, request dto.RegisterRequest) error {
	req := auth.RegisterRequest{
		Email:                request.Email,
		Password:             request.Password,
		PasswordConfirmation: request.PasswordConfirmation,
	}

	if _, err := as.authClient.Register(ctx, &req); err != nil {
		return eris.Wrap(err, appconstant.ErrServiceClient)
	}

	return nil
}

func (as *authServiceGrpc) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error) {
	req := auth.LoginRequest{
		Email:    request.Email,
		Password: request.Password,
	}

	response, err := as.authClient.Login(ctx, &req)
	if err != nil {
		return dto.LoginResponse{}, eris.Wrap(err, appconstant.ErrServiceClient)
	}

	return dto.LoginResponse{
		Type:  response.GetType(),
		Token: response.GetToken(),
	}, nil
}

func (as *authServiceGrpc) VerifyToken(ctx context.Context, token string) (bool, map[string]any, error) {
	data, err := as.authClient.VerifyToken(ctx, &auth.VerifyTokenRequest{
		Token: token,
	})
	if err != nil {
		return false, nil, eris.Wrap(err, appconstant.ErrServiceClient)
	}

	return true, map[string]any{
		appconstant.ContextProfileID: data.GetProfileId(),
	}, nil
}
