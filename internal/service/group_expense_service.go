package service

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/billsplittr/internal/service/fee"
	"github.com/itsLeonB/billsplittr/internal/util"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

type groupExpenseServiceImpl struct {
	transactor                   ezutil.Transactor
	groupExpenseRepository       repository.GroupExpenseRepository
	userService                  UserService
	friendshipService            FriendshipService
	expenseItemRepository        repository.ExpenseItemRepository
	expenseParticipantRepository repository.ExpenseParticipantRepository
	debtService                  DebtService
	feeCalculatorRegistry        map[appconstant.FeeCalculationMethod]fee.FeeCalculator
}

func NewGroupExpenseService(
	transactor ezutil.Transactor,
	groupExpenseRepository repository.GroupExpenseRepository,
	userService UserService,
	friendshipService FriendshipService,
	expenseItemRepository repository.ExpenseItemRepository,
	groupExpenseParticipantRepository repository.ExpenseParticipantRepository,
	debtService DebtService,
) GroupExpenseService {
	return &groupExpenseServiceImpl{
		transactor,
		groupExpenseRepository,
		userService,
		friendshipService,
		expenseItemRepository,
		groupExpenseParticipantRepository,
		debtService,
		fee.NewFeeCalculatorRegistry(),
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

	spec := ezutil.Specification[entity.GroupExpense]{}
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

	spec := ezutil.Specification[entity.GroupExpense]{}
	spec.Model.ID = id
	spec.PreloadRelations = []string{
		"Items",
		"OtherFees",
		"PayerProfile",
		"CreatorProfile",
		"Items.Participants",
		"Items.Participants.Profile",
		"Participants",
		"Participants.Profile",
	}

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

	spec := ezutil.Specification[entity.ExpenseItem]{}
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

	if request.Amount.Cmp(decimal.Zero) <= 0 {
		return dto.ExpenseItemResponse{}, ezutil.UnprocessableEntityError("amount must be more than 0")
	}

	err = ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := ezutil.Specification[entity.ExpenseItem]{}
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

		if ezutil.CompareUUID(request.GroupExpenseID, expenseItem.GroupExpenseID) != 0 {
			return ezutil.UnprocessableEntityError("mismatched group expense ID")
		}

		groupExpenseSpec := ezutil.Specification[entity.GroupExpense]{}
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

func (ges *groupExpenseServiceImpl) ConfirmDraft(ctx context.Context, id uuid.UUID) (dto.GroupExpenseResponse, error) {
	var response dto.GroupExpenseResponse

	err := ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		userProfileID, err := util.GetProfileID(ctx)
		if err != nil {
			return err
		}

		spec := ezutil.Specification[entity.GroupExpense]{}
		spec.Model.ID = id
		spec.PreloadRelations = []string{"Items", "OtherFees", "PayerProfile", "CreatorProfile", "Items.Participants", "Items.Participants.Profile"}
		spec.ForUpdate = true

		groupExpense, err := ges.getGroupExpense(ctx, spec)
		if err != nil {
			return err
		}

		if groupExpense.Confirmed {
			return ezutil.UnprocessableEntityError("already confirmed")
		}

		isAllAnonymous := true
		participantsByProfileID := make(map[uuid.UUID]*entity.ExpenseParticipant)
		for _, item := range groupExpense.Items {
			if len(item.Participants) < 1 {
				return ezutil.UnprocessableEntityError(fmt.Sprintf("item %s does not have participants", item.Name))
			}
			for _, participant := range item.Participants {
				amountToAdd := item.Amount.Mul(participant.Share).Mul(decimal.NewFromInt(int64(item.Quantity)))
				if expenseParticipant, ok := participantsByProfileID[participant.ProfileID]; ok {
					expenseParticipant.ShareAmount = expenseParticipant.ShareAmount.Add(amountToAdd)
				} else {
					expenseParticipant := entity.ExpenseParticipant{
						ParticipantProfileID: participant.ProfileID,
						ShareAmount:          amountToAdd,
					}
					if participant.Profile.IsAnonymous() || participant.ProfileID == userProfileID {
						expenseParticipant.Confirmed = true
					} else {
						isAllAnonymous = false
					}
					participantsByProfileID[participant.ProfileID] = &expenseParticipant
				}
			}
		}

		groupExpenseParticipants := make([]entity.ExpenseParticipant, 0, len(participantsByProfileID))
		for _, expenseParticipant := range participantsByProfileID {
			groupExpenseParticipants = append(groupExpenseParticipants, *expenseParticipant)
		}

		groupExpense.Participants = groupExpenseParticipants
		updatedOtherFees, err := ges.calculateOtherFeeSplits(groupExpense)
		if err != nil {
			return err
		}

		for _, fee := range updatedOtherFees {
			for _, participant := range fee.Participants {
				if expenseParticipant, ok := participantsByProfileID[participant.ProfileID]; !ok {
					return eris.New("missing participant profile from other fee")
				} else {
					expenseParticipant.ShareAmount = expenseParticipant.ShareAmount.Add(participant.ShareAmount)
				}
			}
		}

		updatedGroupExpenseParticipants := make([]entity.ExpenseParticipant, 0, len(participantsByProfileID))
		for _, expenseParticipant := range participantsByProfileID {
			updatedGroupExpenseParticipants = append(updatedGroupExpenseParticipants, *expenseParticipant)
		}

		if err = ges.groupExpenseRepository.SyncParticipants(ctx, groupExpense.ID, updatedGroupExpenseParticipants); err != nil {
			return err
		}

		groupExpense.Confirmed = true
		groupExpense.ParticipantsConfirmed = isAllAnonymous
		// TODO: explore cleaner way
		groupExpense.Participants = nil // Prevent GORM updating child, already synced above
		groupExpense.OtherFees = updatedOtherFees

		updatedGroupExpense, err := ges.groupExpenseRepository.Update(ctx, groupExpense)
		if err != nil {
			return err
		}

		if isAllAnonymous {
			if err = ges.notifyParticipantsConfirmed(ctx, groupExpense.ID); err != nil {
				return err
			}
		} else {
			if err = ges.notifyDraftConfirmed(ctx); err != nil {
				return err
			}
		}

		response = mapper.GroupExpenseToResponse(updatedGroupExpense, userProfileID)

		return nil
	})

	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	return response, nil
}

func (ges *groupExpenseServiceImpl) GetFeeCalculationMethods() []dto.FeeCalculationMethodInfo {
	feeCalculationMethodInfos := make([]dto.FeeCalculationMethodInfo, 0, len(ges.feeCalculatorRegistry))
	for _, feeCalculator := range ges.feeCalculatorRegistry {
		feeCalculationMethodInfos = append(feeCalculationMethodInfos, feeCalculator.GetInfo())
	}

	return feeCalculationMethodInfos
}

func (ges *groupExpenseServiceImpl) calculateOtherFeeSplits(groupExpense entity.GroupExpense) ([]entity.OtherFee, error) {
	var splitErr error

	mapperFunc := func(fee entity.OtherFee) entity.OtherFee {
		feeCalculator, ok := ges.feeCalculatorRegistry[fee.CalculationMethod]
		if !ok {
			splitErr = eris.Errorf("unsupported calculation method: %s", fee.CalculationMethod)
			return entity.OtherFee{}
		}

		if err := feeCalculator.Validate(fee, groupExpense); err != nil {
			splitErr = err
			return entity.OtherFee{}
		}

		fee.Participants = feeCalculator.Split(fee, groupExpense)

		return fee
	}

	splitFees := ezutil.MapSlice(groupExpense.OtherFees, mapperFunc)

	return splitFees, splitErr
}

func (ges *groupExpenseServiceImpl) notifyParticipantsConfirmed(ctx context.Context, groupExpenseID uuid.UUID) error {
	if os.Getenv("ENABLE_ASYNC") == "true" {
		panic("to be implemented")
	}

	return ges.debtService.ProcessConfirmedGroupExpense(ctx, groupExpenseID)
}

func (ges *groupExpenseServiceImpl) notifyDraftConfirmed(ctx context.Context) error {
	panic("to be implemented")
}

func (ges *groupExpenseServiceImpl) getGroupExpense(ctx context.Context, spec ezutil.Specification[entity.GroupExpense]) (entity.GroupExpense, error) {
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
