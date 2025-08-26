package service

import (
	"context"
	"log"

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
		if isFriend, _, err := ebs.friendshipService.IsFriends(ctx, request.CreatorProfileID, request.PayerProfileID); err != nil {
			return err
		} else if !isFriend {
			return ezutil.UnprocessableEntityError(appconstant.ErrNotFriends)
		}
	}

	defer func() {
		if err := request.ImageReader.Close(); err != nil {
			log.Printf("error closing ImageReader: %v\n", err)
		}
	}()

	filename, err := ebs.imageRepository.Upload(ctx, request.ImageReader, request.ContentType)
	if err != nil {
		return err
	}

	expenseBill := entity.ExpenseBill{
		PayerProfileID:   request.PayerProfileID,
		ImageName:        filename,
		CreatorProfileID: request.CreatorProfileID,
	}

	if _, err = ebs.expenseBillRepository.Insert(ctx, expenseBill); err != nil {
		if err = ebs.imageRepository.Delete(ctx, filename); err != nil {
			log.Printf("warning: failed to clean up GCS image '%s': %v", filename, err)
		}
		return err
	}

	return nil
}
