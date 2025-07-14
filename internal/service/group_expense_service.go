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
	transactor             ezutil.Transactor
	groupExpenseRepository repository.GroupExpenseRepository
	userService            UserService
	friendshipService      FriendshipService
	expenseItemRepository  repository.ExpenseItemRepository
}

func NewGroupExpenseService(
	transactor ezutil.Transactor,
	groupExpenseRepository repository.GroupExpenseRepository,
	userService UserService,
	friendshipService FriendshipService,
	expenseItemRepository repository.ExpenseItemRepository,
) GroupExpenseService {
	return &groupExpenseServiceImpl{
		transactor,
		groupExpenseRepository,
		userService,
		friendshipService,
		expenseItemRepository,
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

	spec := entity.GenericSpec[entity.GroupExpense]{}
	spec.Model.CreatorProfileID = user.Profile.ID
	spec.PreloadRelations = []string{"Items", "OtherFees", "PayerProfile", "CreatorProfile"}

	groupExpenses, err := ges.groupExpenseRepository.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(groupExpenses, mapper.GetGroupExpenseSimpleMapper(user.Profile.ID)), nil
}

func (ges *groupExpenseServiceImpl) GetDetails(ctx context.Context, id uuid.UUID) (dto.GroupExpenseResponse, error) {
	profileID, err := util.GetProfileID(ctx)
	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	spec := entity.GenericSpec[entity.GroupExpense]{}
	spec.Model.ID = id
	spec.PreloadRelations = []string{"Items", "OtherFees", "PayerProfile", "CreatorProfile"}

	groupExpense, err := ges.getGroupExpense(ctx, spec)
	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	return mapper.GroupExpenseToResponse(groupExpense, profileID), nil
}

func (ges *groupExpenseServiceImpl) GetItemDetails(ctx context.Context, groupExpenseID, expenseItemID uuid.UUID) (dto.ExpenseItemResponse, error) {
	profileID, err := util.GetProfileID(ctx)
	if err != nil {
		return dto.ExpenseItemResponse{}, err
	}

	spec := entity.GenericSpec[entity.ExpenseItem]{}
	spec.Model.ID = expenseItemID
	spec.Model.GroupExpenseID = groupExpenseID
	spec.PreloadRelations = []string{"Participants", "Participants.Profile"}

	expenseItem, err := ges.expenseItemRepository.FindFirst(ctx, spec)
	if err != nil {
		return dto.ExpenseItemResponse{}, err
	}
	if expenseItem.IsZero() {
		return dto.ExpenseItemResponse{}, ezutil.NotFoundError(util.NotFoundMessage(spec.Model))
	}
	if expenseItem.IsDeleted() {
		return dto.ExpenseItemResponse{}, ezutil.UnprocessableEntityError(util.DeletedMessage(expenseItem))
	}

	return mapper.ExpenseItemToResponse(expenseItem, profileID), nil
}

func (ges *groupExpenseServiceImpl) UpdateItem(ctx context.Context, request dto.UpdateExpenseItemRequest) (dto.ExpenseItemResponse, error) {
	var response dto.ExpenseItemResponse

	profileID, err := util.GetProfileID(ctx)
	if err != nil {
		return dto.ExpenseItemResponse{}, err
	}

	if request.Amount.Cmp(decimal.Zero) == 0 {
		return dto.ExpenseItemResponse{}, ezutil.UnprocessableEntityError("amount must be more than 0")
	}

	err = ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := entity.GenericSpec[entity.ExpenseItem]{}
		spec.Model.ID = request.ID
		spec.ForUpdate = true

		expenseItem, err := ges.expenseItemRepository.FindFirst(ctx, spec)
		if err != nil {
			return err
		}
		if expenseItem.IsZero() {
			return ezutil.NotFoundError(util.NotFoundMessage(spec.Model))
		}
		if expenseItem.IsDeleted() {
			return ezutil.UnprocessableEntityError(util.DeletedMessage(expenseItem))
		}

		if util.CompareUUID(request.GroupExpenseID, expenseItem.GroupExpenseID) != 0 {
			return ezutil.UnprocessableEntityError("mismatched group expense ID")
		}

		groupExpenseSpec := entity.GenericSpec[entity.GroupExpense]{}
		groupExpenseSpec.Model.ID = expenseItem.GroupExpenseID
		groupExpenseSpec.ForUpdate = true

		groupExpense, err := ges.getGroupExpense(ctx, groupExpenseSpec)
		if err != nil {
			return err
		}

		patchedExpenseItem := mapper.PatchExpenseItemWithRequest(expenseItem, request)

		updatedExpenseItem, err := ges.expenseItemRepository.Update(ctx, patchedExpenseItem)
		if err != nil {
			return err
		}

		oldAmount := expenseItem.Amount.Mul(decimal.NewFromInt(int64(expenseItem.Quantity)))
		newAmount := updatedExpenseItem.Amount.Mul(decimal.NewFromInt(int64(updatedExpenseItem.Quantity)))

		groupExpense.TotalAmount = groupExpense.TotalAmount.
			Sub(oldAmount).
			Add(newAmount)

		if _, err := ges.groupExpenseRepository.Update(ctx, groupExpense); err != nil {
			return err
		}

		updatedParticipants := ezutil.MapSlice(request.Participants, mapper.ItemParticipantRequestToEntity)
		if err := ges.expenseItemRepository.SyncParticipants(ctx, updatedExpenseItem.ID, updatedParticipants); err != nil {
			return err
		}

		response = mapper.ExpenseItemToResponse(updatedExpenseItem, profileID)

		return nil
	})

	if err != nil {
		return dto.ExpenseItemResponse{}, err
	}

	return response, nil
}

func (ges *groupExpenseServiceImpl) getGroupExpense(ctx context.Context, spec entity.GenericSpec[entity.GroupExpense]) (entity.GroupExpense, error) {
	groupExpense, err := ges.groupExpenseRepository.FindFirst(ctx, spec)
	if err != nil {
		return entity.GroupExpense{}, err
	}
	if groupExpense.IsZero() {
		return entity.GroupExpense{}, ezutil.NotFoundError(util.NotFoundMessage(spec.Model))
	}
	if groupExpense.IsDeleted() {
		return entity.GroupExpense{}, ezutil.UnprocessableEntityError(util.DeletedMessage(groupExpense))
	}

	return groupExpense, nil
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
