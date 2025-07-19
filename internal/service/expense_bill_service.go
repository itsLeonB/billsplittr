package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/ezutil"
)

type expenseBillServiceImpl struct {
	friendshipService     FriendshipService
	imageRepository       repository.ImageRepository
	expenseBillRepository repository.ExpenseBillRepository
}

func NewExpenseBillService(
	friendshipService FriendshipService,
	imageRepository repository.ImageRepository,
	expenseBillRepository repository.ExpenseBillRepository,
) ExpenseBillService {
	return &expenseBillServiceImpl{
		friendshipService,
		imageRepository,
		expenseBillRepository,
	}
}

func (ebs *expenseBillServiceImpl) Upload(ctx context.Context, request dto.NewExpenseBillRequest) error {
	// Default PayerProfileID to the user's profile ID if not provided
	// This is useful when the user is creating a group expense for themselves.
	if request.PayerProfileID == uuid.Nil {
		request.PayerProfileID = request.CreatorProfileID
	} else {
		// Check if the payer is a friend of the user
		if isFriend, err := ebs.friendshipService.IsFriends(ctx, request.CreatorProfileID, request.PayerProfileID); err != nil {
			return err
		} else if !isFriend {
			return ezutil.UnprocessableEntityError(appconstant.ErrNotFriends)
		}
	}

	filename, err := ebs.imageRepository.Upload(ctx, request.ImageFile)
	if err != nil {
		return err
	}

	expenseBill := entity.ExpenseBill{
		PayerProfileID:   request.PayerProfileID,
		ImageName:        filename,
		CreatorProfileID: request.CreatorProfileID,
	}

	_, err = ebs.expenseBillRepository.Insert(ctx, expenseBill)

	return err
}
