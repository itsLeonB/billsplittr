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
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

type groupExpenseServiceImpl struct {
	transactor             crud.Transactor
	groupExpenseRepository repository.GroupExpenseRepository
	friendshipService      FriendshipService
	expenseItemRepository  repository.ExpenseItemRepository
	debtService            DebtService
	feeCalculatorRegistry  map[appconstant.FeeCalculationMethod]fee.FeeCalculator
	otherFeeRepository     repository.OtherFeeRepository
	profileService         ProfileService
}

func NewGroupExpenseService(
	transactor crud.Transactor,
	groupExpenseRepository repository.GroupExpenseRepository,
	friendshipService FriendshipService,
	expenseItemRepository repository.ExpenseItemRepository,
	debtService DebtService,
	otherFeeRepository repository.OtherFeeRepository,
	profileService ProfileService,
) GroupExpenseService {
	return &groupExpenseServiceImpl{
		transactor,
		groupExpenseRepository,
		friendshipService,
		expenseItemRepository,
		debtService,
		fee.NewFeeCalculatorRegistry(),
		otherFeeRepository,
		profileService,
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

	namesByProfileIDs, err := ges.getGroupExpenseProfileNames(ctx, insertedGroupExpense)
	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	return mapper.GroupExpenseToResponse(insertedGroupExpense, request.CreatedByProfileID, namesByProfileIDs), nil
}

func (ges *groupExpenseServiceImpl) GetAllCreated(ctx context.Context, profileID uuid.UUID) ([]dto.GroupExpenseResponse, error) {
	spec := crud.Specification[entity.GroupExpense]{}
	spec.Model.CreatorProfileID = profileID
	spec.PreloadRelations = []string{"Items", "OtherFees"}

	groupExpenses, err := ges.groupExpenseRepository.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	profileIDs := make([]uuid.UUID, 0)
	for _, groupExpense := range groupExpenses {
		profileIDs = append(profileIDs, groupExpense.ProfileIDs()...)
	}

	namesByProfileIDs := make(map[uuid.UUID]string, len(profileIDs))
	if len(profileIDs) > 0 {
		namesByProfileIDs, err = ges.profileService.GetNames(ctx, profileIDs)
		if err != nil {
			return nil, err
		}
	}

	mapFunc := func(groupExpense entity.GroupExpense) dto.GroupExpenseResponse {
		return mapper.GroupExpenseToResponse(groupExpense, profileID, namesByProfileIDs)
	}

	return ezutil.MapSlice(groupExpenses, mapFunc), nil
}

func (ges *groupExpenseServiceImpl) GetDetails(ctx context.Context, id uuid.UUID, profileID uuid.UUID) (dto.GroupExpenseResponse, error) {
	spec := crud.Specification[entity.GroupExpense]{}
	spec.Model.ID = id
	spec.PreloadRelations = []string{
		"Items",
		"OtherFees",
		"Items.Participants",
		"Participants",
	}

	groupExpense, err := ges.getGroupExpense(ctx, spec)
	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	namesByProfileIDs, err := ges.getGroupExpenseProfileNames(ctx, groupExpense)
	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	return mapper.GroupExpenseToResponse(groupExpense, profileID, namesByProfileIDs), nil
}

func (ges *groupExpenseServiceImpl) GetItemDetails(ctx context.Context, groupExpenseID, expenseItemID, profileID uuid.UUID) (dto.ExpenseItemResponse, error) {
	spec := crud.Specification[entity.ExpenseItem]{}
	spec.Model.ID = expenseItemID
	spec.Model.GroupExpenseID = groupExpenseID
	spec.PreloadRelations = []string{"Participants"}

	expenseItem, err := ges.getExpenseItemBySpec(ctx, spec)
	if err != nil {
		return dto.ExpenseItemResponse{}, err
	}

	namesByProfileIDs, err := ges.profileService.GetNames(ctx, expenseItem.ProfileIDs())
	if err != nil {
		return dto.ExpenseItemResponse{}, err
	}

	return mapper.ExpenseItemToResponse(expenseItem, profileID, namesByProfileIDs), nil
}

func (ges *groupExpenseServiceImpl) UpdateItem(ctx context.Context, profileID uuid.UUID, request dto.UpdateExpenseItemRequest) (dto.ExpenseItemResponse, error) {
	var response dto.ExpenseItemResponse

	if !request.Amount.IsPositive() {
		return dto.ExpenseItemResponse{}, ungerr.UnprocessableEntityError(appconstant.ErrNonPositiveAmount)
	}

	for _, participant := range request.Participants {
		if participant.ProfileID == profileID {
			continue // skip if participant is current user
		}
		// Check if the participant is a friend of the user
		if isFriend, _, err := ges.friendshipService.IsFriends(ctx, profileID, participant.ProfileID); err != nil {
			return dto.ExpenseItemResponse{}, err
		} else if !isFriend {
			return dto.ExpenseItemResponse{}, ungerr.UnprocessableEntityError(appconstant.ErrNotFriends)
		}
	}

	err := ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		expenseItem, err := ges.getExpenseItemByIDForUpdate(ctx, request.ID, request.GroupExpenseID)
		if err != nil {
			return err
		}

		if ezutil.CompareUUID(request.GroupExpenseID, expenseItem.GroupExpenseID) != 0 {
			return ungerr.UnprocessableEntityError("mismatched group expense ID")
		}

		groupExpense, err := ges.getUnconfirmedGroupExpenseForUpdate(ctx, expenseItem.GroupExpenseID)
		if err != nil {
			return err
		}

		patchedExpenseItem := mapper.PatchExpenseItemWithRequest(expenseItem, request)

		updatedExpenseItem, err := ges.expenseItemRepository.Update(ctx, patchedExpenseItem)
		if err != nil {
			return err
		}

		oldAmount := expenseItem.TotalAmount()
		newAmount := updatedExpenseItem.TotalAmount()

		if oldAmount.Cmp(newAmount) != 0 {
			groupExpense.TotalAmount = groupExpense.TotalAmount.
				Sub(oldAmount).
				Add(newAmount)

			if _, err := ges.groupExpenseRepository.Update(ctx, groupExpense); err != nil {
				return err
			}
		}

		updatedParticipants := ezutil.MapSlice(request.Participants, mapper.ItemParticipantRequestToEntity)
		if err := ges.expenseItemRepository.SyncParticipants(ctx, updatedExpenseItem.ID, updatedParticipants); err != nil {
			return err
		}

		updatedExpenseItem.Participants = updatedParticipants
		profileIDs := updatedExpenseItem.ProfileIDs()
		namesByProfileIDs, err := ges.profileService.GetNames(ctx, profileIDs)
		if err != nil {
			return err
		}

		response = mapper.ExpenseItemToResponse(updatedExpenseItem, profileID, namesByProfileIDs)

		return nil
	})

	if err != nil {
		return dto.ExpenseItemResponse{}, err
	}

	return response, nil
}

func (ges *groupExpenseServiceImpl) ConfirmDraft(ctx context.Context, id, profileID uuid.UUID) (dto.GroupExpenseResponse, error) {
	var response dto.GroupExpenseResponse

	err := ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[entity.GroupExpense]{}
		spec.Model.ID = id
		spec.PreloadRelations = []string{"Items", "OtherFees", "Items.Participants"}
		spec.ForUpdate = true

		groupExpense, err := ges.getGroupExpense(ctx, spec)
		if err != nil {
			return err
		}

		if groupExpense.Confirmed {
			return ungerr.UnprocessableEntityError("already confirmed")
		}

		profileIDs := groupExpense.ProfileIDs()
		profilesByID, err := ges.profileService.GetByIDs(ctx, profileIDs)
		if err != nil {
			return err
		}

		isAllAnonymous := true
		participantsByProfileID := make(map[uuid.UUID]*entity.ExpenseParticipant)
		for _, item := range groupExpense.Items {
			if len(item.Participants) < 1 {
				return ungerr.UnprocessableEntityError(fmt.Sprintf("item %s does not have participants", item.Name))
			}
			for _, participant := range item.Participants {
				amountToAdd := item.TotalAmount().Mul(participant.Share)
				if expenseParticipant, ok := participantsByProfileID[participant.ProfileID]; ok {
					expenseParticipant.ShareAmount = expenseParticipant.ShareAmount.Add(amountToAdd)
				} else {
					expenseParticipant := entity.ExpenseParticipant{
						ParticipantProfileID: participant.ProfileID,
						ShareAmount:          amountToAdd,
					}
					if profilesByID[participant.ProfileID].IsAnonymous || participant.ProfileID == profileID {
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
		updatedOtherFees, err := ges.calculateOtherFeeSplits(ctx, groupExpense)
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

		updatedGroupExpense, err := ges.groupExpenseRepository.Update(ctx, groupExpense)
		if err != nil {
			return err
		}

		updatedGroupExpense.Participants = updatedGroupExpenseParticipants

		if isAllAnonymous {
			if err = ges.notifyParticipantsConfirmed(ctx, updatedGroupExpense); err != nil {
				return err
			}
		} else {
			if err = ges.notifyDraftConfirmed(ctx); err != nil {
				return err
			}
		}

		namesByProfileID, err := ges.getGroupExpenseProfileNames(ctx, updatedGroupExpense)
		if err != nil {
			return err
		}

		response = mapper.GroupExpenseToResponse(updatedGroupExpense, profileID, namesByProfileID)

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

func (ges *groupExpenseServiceImpl) UpdateFee(ctx context.Context, profileID uuid.UUID, request dto.UpdateOtherFeeRequest) (dto.OtherFeeResponse, error) {
	var response dto.OtherFeeResponse

	if request.Amount.Cmp(decimal.Zero) <= 0 {
		return dto.OtherFeeResponse{}, ungerr.UnprocessableEntityError("amount must be more than 0")
	}

	err := ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		groupExpense, err := ges.getUnconfirmedGroupExpenseForUpdate(ctx, request.GroupExpenseID)
		if err != nil {
			return err
		}

		spec := crud.Specification[entity.OtherFee]{}
		spec.Model.ID = request.ID
		spec.Model.GroupExpenseID = request.GroupExpenseID
		spec.ForUpdate = true
		otherFee, err := ges.otherFeeRepository.FindFirst(ctx, spec)
		if err != nil {
			return err
		}
		if otherFee.IsZero() {
			return ungerr.NotFoundError(fmt.Sprintf("other fee with ID: %s is not found", request.ID))
		}
		if otherFee.IsDeleted() {
			return ungerr.UnprocessableEntityError(fmt.Sprintf("other fee with ID: %s is deleted", request.ID))
		}

		patchedFee := mapper.PatchOtherFeeWithRequest(otherFee, request)

		updatedFee, err := ges.otherFeeRepository.Update(ctx, patchedFee)
		if err != nil {
			return err
		}

		if updatedFee.Amount.Cmp(otherFee.Amount) != 0 {
			groupExpense.TotalAmount = groupExpense.TotalAmount.Sub(otherFee.Amount).Add(updatedFee.Amount)
			if _, err = ges.groupExpenseRepository.Update(ctx, groupExpense); err != nil {
				return err
			}
		}

		profileIDs := updatedFee.ProfileIDs()
		namesByProfileIDs, err := ges.profileService.GetNames(ctx, profileIDs)
		if err != nil {
			return err
		}

		response = mapper.OtherFeeToResponse(updatedFee, profileID, namesByProfileIDs)

		return nil
	})

	if err != nil {
		return dto.OtherFeeResponse{}, err
	}

	return response, nil
}

func (ges *groupExpenseServiceImpl) AddItem(ctx context.Context, profileID uuid.UUID, request dto.NewExpenseItemRequest) (dto.ExpenseItemResponse, error) {
	var response dto.ExpenseItemResponse

	if !request.Amount.IsPositive() {
		return dto.ExpenseItemResponse{}, ungerr.UnprocessableEntityError(appconstant.ErrNonPositiveAmount)
	}

	err := ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		groupExpense, err := ges.getUnconfirmedGroupExpenseForUpdate(ctx, request.GroupExpenseID)
		if err != nil {
			return err
		}

		expenseItem := mapper.ExpenseItemRequestToEntity(request)

		itemTotalAmount := expenseItem.TotalAmount()
		groupExpense.TotalAmount = groupExpense.TotalAmount.Add(itemTotalAmount)
		groupExpense.Subtotal = groupExpense.Subtotal.Add(itemTotalAmount)
		if _, err = ges.groupExpenseRepository.Update(ctx, groupExpense); err != nil {
			return err
		}

		insertedItem, err := ges.expenseItemRepository.Insert(ctx, expenseItem)
		if err != nil {
			return err
		}

		profileIDs := insertedItem.ProfileIDs()
		namesByProfileID, err := ges.profileService.GetNames(ctx, profileIDs)
		if err != nil {
			return err
		}

		response = mapper.ExpenseItemToResponse(insertedItem, profileID, namesByProfileID)

		return nil
	})

	if err != nil {
		return dto.ExpenseItemResponse{}, err
	}

	return response, nil
}

func (ges *groupExpenseServiceImpl) AddFee(ctx context.Context, profileID uuid.UUID, request dto.NewOtherFeeRequest) (dto.OtherFeeResponse, error) {
	var response dto.OtherFeeResponse

	if !request.Amount.IsPositive() {
		return dto.OtherFeeResponse{}, ungerr.UnprocessableEntityError(appconstant.ErrNonPositiveAmount)
	}

	err := ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		groupExpense, err := ges.getUnconfirmedGroupExpenseForUpdate(ctx, request.GroupExpenseID)
		if err != nil {
			return err
		}

		fee := mapper.OtherFeeRequestToEntity(request)

		groupExpense.TotalAmount = groupExpense.TotalAmount.Add(fee.Amount)
		if _, err = ges.groupExpenseRepository.Update(ctx, groupExpense); err != nil {
			return err
		}

		insertedFee, err := ges.otherFeeRepository.Insert(ctx, fee)
		if err != nil {
			return err
		}

		profileIDs := insertedFee.ProfileIDs()
		namesByProfileID, err := ges.profileService.GetNames(ctx, profileIDs)
		if err != nil {
			return err
		}

		response = mapper.OtherFeeToResponse(insertedFee, profileID, namesByProfileID)

		return nil
	})

	if err != nil {
		return dto.OtherFeeResponse{}, err
	}

	return response, nil
}

func (ges *groupExpenseServiceImpl) RemoveItem(ctx context.Context, request dto.DeleteExpenseItemRequest) error {
	return ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		groupExpense, err := ges.getUnconfirmedGroupExpenseForUpdate(ctx, request.GroupExpenseID)
		if err != nil {
			return err
		}

		expenseItem, err := ges.getExpenseItemByIDForUpdate(ctx, request.ID, request.GroupExpenseID)
		if err != nil {
			return err
		}

		if err = ges.expenseItemRepository.Delete(ctx, expenseItem); err != nil {
			return err
		}

		itemAmount := expenseItem.TotalAmount()
		groupExpense.TotalAmount = groupExpense.TotalAmount.Sub(itemAmount)
		groupExpense.Subtotal = groupExpense.Subtotal.Sub(itemAmount)

		if _, err = ges.groupExpenseRepository.Update(ctx, groupExpense); err != nil {
			return err
		}

		return nil
	})
}

func (ges *groupExpenseServiceImpl) RemoveFee(ctx context.Context, request dto.DeleteOtherFeeRequest) error {
	return ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		groupExpense, err := ges.getUnconfirmedGroupExpenseForUpdate(ctx, request.GroupExpenseID)
		if err != nil {
			return err
		}

		spec := crud.Specification[entity.OtherFee]{}
		spec.Model.ID = request.ID
		spec.Model.GroupExpenseID = request.GroupExpenseID
		spec.ForUpdate = true
		otherFee, err := ges.otherFeeRepository.FindFirst(ctx, spec)
		if err != nil {
			return err
		}
		if otherFee.IsZero() {
			return ungerr.NotFoundError(fmt.Sprintf("other fee with ID: %s is not found", request.ID))
		}
		if otherFee.IsDeleted() {
			return ungerr.UnprocessableEntityError(fmt.Sprintf("other fee with ID: %s is deleted", request.ID))
		}

		if err = ges.otherFeeRepository.Delete(ctx, otherFee); err != nil {
			return err
		}

		groupExpense.TotalAmount = groupExpense.TotalAmount.Sub(otherFee.Amount)
		if _, err = ges.groupExpenseRepository.Update(ctx, groupExpense); err != nil {
			return err
		}

		return nil
	})
}

func (ges *groupExpenseServiceImpl) getGroupExpenseProfileNames(ctx context.Context, groupExpense entity.GroupExpense) (map[uuid.UUID]string, error) {
	namesByProfileIDs, err := ges.profileService.GetNames(ctx, groupExpense.ProfileIDs())
	if err != nil {
		return nil, err
	}

	return namesByProfileIDs, nil
}

func (ges *groupExpenseServiceImpl) getExpenseItemByIDForUpdate(ctx context.Context, expenseItemID, groupExpenseID uuid.UUID) (entity.ExpenseItem, error) {
	spec := crud.Specification[entity.ExpenseItem]{}
	spec.Model.ID = expenseItemID
	spec.Model.GroupExpenseID = groupExpenseID
	spec.ForUpdate = true

	expenseItem, err := ges.getExpenseItemBySpec(ctx, spec)
	if err != nil {
		return entity.ExpenseItem{}, err
	}

	return expenseItem, nil
}

func (ges *groupExpenseServiceImpl) calculateOtherFeeSplits(ctx context.Context, groupExpense entity.GroupExpense) ([]entity.OtherFee, error) {
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

		if err := ges.otherFeeRepository.SyncParticipants(ctx, fee.ID, fee.Participants); err != nil {
			splitErr = err
			return entity.OtherFee{}
		}

		return fee
	}

	splitFees := ezutil.MapSlice(groupExpense.OtherFees, mapperFunc)

	return splitFees, splitErr
}

func (ges *groupExpenseServiceImpl) notifyParticipantsConfirmed(ctx context.Context, groupExpense entity.GroupExpense) error {
	if os.Getenv("ENABLE_ASYNC") == "true" {
		panic("to be implemented")
	}

	return ges.debtService.ProcessConfirmedGroupExpense(ctx, groupExpense)
}

func (ges *groupExpenseServiceImpl) notifyDraftConfirmed(ctx context.Context) error {
	panic("to be implemented")
}

func (ges *groupExpenseServiceImpl) getUnconfirmedGroupExpenseForUpdate(ctx context.Context, id uuid.UUID) (entity.GroupExpense, error) {
	spec := crud.Specification[entity.GroupExpense]{}
	spec.Model.ID = id
	spec.ForUpdate = true
	groupExpense, err := ges.getGroupExpense(ctx, spec)
	if err != nil {
		return entity.GroupExpense{}, err
	}
	if groupExpense.Confirmed || groupExpense.ParticipantsConfirmed {
		return entity.GroupExpense{}, ungerr.UnprocessableEntityError("expense already confirmed")
	}

	return groupExpense, nil
}

func (ges *groupExpenseServiceImpl) getGroupExpense(ctx context.Context, spec crud.Specification[entity.GroupExpense]) (entity.GroupExpense, error) {
	groupExpense, err := ges.groupExpenseRepository.FindFirst(ctx, spec)
	if err != nil {
		return entity.GroupExpense{}, err
	}
	if groupExpense.IsZero() {
		return entity.GroupExpense{}, ungerr.NotFoundError(fmt.Sprintf("group expense with ID %s is not found", spec.Model.ID))
	}
	if groupExpense.IsDeleted() {
		return entity.GroupExpense{}, ungerr.UnprocessableEntityError(fmt.Sprintf("group expense with ID %s is deleted", spec.Model.ID))
	}

	return groupExpense, nil
}

func (ges *groupExpenseServiceImpl) validateAndPatchRequest(ctx context.Context, request *dto.NewGroupExpenseRequest) error {
	if request.TotalAmount.IsZero() {
		return ungerr.UnprocessableEntityError(appconstant.ErrAmountZero)
	}

	calculatedFeeTotal := decimal.Zero
	calculatedSubtotal := decimal.Zero
	for _, item := range request.Items {
		calculatedSubtotal = calculatedSubtotal.Add(item.Amount.Mul(decimal.NewFromInt(int64(item.Quantity))))
	}
	for _, fee := range request.OtherFees {
		calculatedFeeTotal = calculatedFeeTotal.Add(fee.Amount)
	}
	if calculatedFeeTotal.Add(calculatedSubtotal).Cmp(request.TotalAmount) != 0 {
		return ungerr.UnprocessableEntityError(appconstant.ErrAmountMismatched)
	}
	if calculatedSubtotal.Cmp(request.Subtotal) != 0 {
		return ungerr.UnprocessableEntityError(appconstant.ErrAmountMismatched)
	}

	// Default PayerProfileID to the user's profile ID if not provided
	// This is useful when the user is creating a group expense for themselves.
	if request.PayerProfileID == uuid.Nil {
		request.PayerProfileID = request.CreatedByProfileID
	} else {
		// Check if the payer is a friend of the user
		isFriend, _, err := ges.friendshipService.IsFriends(ctx, request.CreatedByProfileID, request.PayerProfileID)
		if err != nil {
			return err
		}
		if !isFriend {
			return ungerr.UnprocessableEntityError(appconstant.ErrNotFriends)
		}
	}

	return nil
}

func (ges *groupExpenseServiceImpl) getExpenseItemBySpec(ctx context.Context, spec crud.Specification[entity.ExpenseItem]) (entity.ExpenseItem, error) {
	expenseItem, err := ges.expenseItemRepository.FindFirst(ctx, spec)
	if err != nil {
		return entity.ExpenseItem{}, err
	}
	if expenseItem.IsZero() {
		return entity.ExpenseItem{}, ungerr.NotFoundError(fmt.Sprintf("expense item with ID %s is not found", spec.Model.ID))
	}
	if expenseItem.IsDeleted() {
		return entity.ExpenseItem{}, ungerr.UnprocessableEntityError(fmt.Sprintf("expense item with ID %s is deleted", spec.Model.ID))
	}

	return expenseItem, nil
}
