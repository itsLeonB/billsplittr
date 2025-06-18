package service

import (
	"context"
	"fmt"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/ezutil"
)

type authServiceImpl struct {
	hashService    ezutil.HashService
	jwtService     ezutil.JWTService
	userRepository repository.UserRepository
	transactor     ezutil.Transactor
}

func NewAuthService(
	hashService ezutil.HashService,
	jwtService ezutil.JWTService,
	userRepository repository.UserRepository,
	transactor ezutil.Transactor,
) AuthService {
	return &authServiceImpl{
		hashService:    hashService,
		jwtService:     jwtService,
		userRepository: userRepository,
		transactor:     transactor,
	}
}

func (as *authServiceImpl) Register(ctx context.Context, request dto.RegisterRequest) error {
	return ezutil.WithinTransaction(ctx, as.transactor, func(ctx context.Context) error {
		spec := entity.User{Email: request.Email}

		existingUser, err := as.userRepository.Find(ctx, spec)
		if err != nil {
			return err
		}
		if !existingUser.IsZero() {
			return ezutil.ConflictError(fmt.Sprintf(appconstant.MsgAuthDuplicateUser, request.Email))
		}

		hash, err := as.hashService.Hash(request.Password)
		if err != nil {
			return err
		}

		spec.Password = hash

		_, err = as.userRepository.Insert(ctx, spec)
		if err != nil {
			return err
		}

		return nil
	})
}

func (as *authServiceImpl) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error) {
	spec := entity.User{Email: request.Email}

	user, err := as.userRepository.Find(ctx, spec)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if user.IsZero() {
		return dto.LoginResponse{}, ezutil.NotFoundError(appconstant.MsgAuthUnknownCredentials)
	}

	ok, err := as.hashService.CheckHash(user.Password, request.Password)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if !ok {
		return dto.LoginResponse{}, ezutil.NotFoundError(appconstant.MsgAuthUnknownCredentials)
	}

	token, err := as.jwtService.CreateToken(mapper.UserToAuthData(user))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		Type:  "Bearer",
		Token: token,
	}, nil
}
