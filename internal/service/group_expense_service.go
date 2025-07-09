package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/billsplittr/internal/util"
	"github.com/itsLeonB/ezutil"
	"github.com/shopspring/decimal"
)

type groupExpenseServiceImpl struct {
	groupExpenseRepository repository.GroupExpenseRepository
	userService            UserService
	friendshipService      FriendshipService
}

func NewGroupExpenseService(
	groupExpenseRepository repository.GroupExpenseRepository,
	userService UserService,
	friendshipService FriendshipService,
) GroupExpenseService {
	return &groupExpenseServiceImpl{
		groupExpenseRepository,
		userService,
		friendshipService,
	}
}

func (ges *groupExpenseServiceImpl) CreateDraft(ctx context.Context, request dto.NewGroupExpenseRequest) (dto.GroupExpenseResponse, error) {
	if err := ges.validateAndPatchRequest(ctx, &request); err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	groupExpense := mapper.GroupExpenseRequestToEntity(request)

	insertedGroupExpense, err := ges.groupExpenseRepository.Insert(ctx, groupExpense)
	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	return mapper.GroupExpenseToResponse(insertedGroupExpense, uuid.Nil), nil
}

func (ges *groupExpenseServiceImpl) GetAllCreated(ctx context.Context, userID uuid.UUID) ([]dto.GroupExpenseResponse, error) {
	user, err := ges.userService.GetEntityByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	spec := entity.GroupExpenseSpecification{}
	spec.CreatorProfileID = user.Profile.ID
	spec.PreloadRelations = []string{"Items", "OtherFees", "PayerProfile", "CreatorProfile"}

	groupExpenses, err := ges.groupExpenseRepository.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(groupExpenses, mapper.GetGroupExpenseSimpleMapper(user.Profile.ID)), nil
}

func (ges *groupExpenseServiceImpl) GetDetails(ctx context.Context, id uuid.UUID) (dto.GroupExpenseResponse, error) {
	userID, err := util.FindUserID(ctx)
	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	user, err := ges.userService.GetEntityByID(ctx, userID)
	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	spec := entity.GroupExpenseSpecification{}
	spec.ID = id
	spec.PreloadRelations = []string{"Items", "OtherFees", "PayerProfile", "CreatorProfile"}

	groupExpense, err := ges.groupExpenseRepository.FindFirst(ctx, spec)
	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}
	if groupExpense.IsZero() {
		return dto.GroupExpenseResponse{}, ezutil.NotFoundError(appconstant.ErrGroupExpenseNotFound(id))
	}
	if groupExpense.IsDeleted() {
		return dto.GroupExpenseResponse{}, ezutil.NotFoundError(appconstant.ErrGroupExpenseDeleted(id))
	}

	return mapper.GroupExpenseToResponse(groupExpense, user.Profile.ID), nil
}

func (ges *groupExpenseServiceImpl) validateAndPatchRequest(ctx context.Context, request *dto.NewGroupExpenseRequest) error {
	if request.TotalAmount.IsZero() {
		return ezutil.UnprocessableEntityError(appconstant.ErrAmountZero)
	}

	calculatedTotal := decimal.Zero
	for _, item := range request.Items {
		calculatedTotal = calculatedTotal.Add(item.Amount.Mul(decimal.NewFromInt(int64(item.Quantity))))
	}
	for _, fee := range request.OtherFees {
		calculatedTotal = calculatedTotal.Add(fee.Amount)
	}
	if calculatedTotal.Cmp(request.TotalAmount) != 0 {
		return ezutil.UnprocessableEntityError(appconstant.ErrAmountMismatched)
	}

	user, err := ges.userService.GetEntityByID(ctx, request.CreatedByUserID)
	if err != nil {
		return err
	}

	request.CreatedByProfileID = user.Profile.ID

	// Default PayerProfileID to the user's profile ID if not provided
	// This is useful when the user is creating a group expense for themselves.
	if request.PayerProfileID == uuid.Nil {
		request.PayerProfileID = user.Profile.ID
	} else {
		// Check if the payer is a friend of the user
		isFriend, err := ges.friendshipService.IsFriends(ctx, user.Profile.ID, request.PayerProfileID)
		if err != nil {
			return err
		}
		if !isFriend {
			return ezutil.UnprocessableEntityError(appconstant.ErrNotFriends)
		}
	}

	return nil
}
