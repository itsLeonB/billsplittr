package service

import (
	"context"
	"fmt"

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
	feeCalculatorRegistry  map[appconstant.FeeCalculationMethod]fee.FeeCalculator
	otherFeeRepository     repository.OtherFeeRepository
}

func NewGroupExpenseService(
	transactor crud.Transactor,
	groupExpenseRepository repository.GroupExpenseRepository,
	otherFeeRepository repository.OtherFeeRepository,
) GroupExpenseService {
	return &groupExpenseServiceImpl{
		transactor,
		groupExpenseRepository,
		fee.NewFeeCalculatorRegistry(),
		otherFeeRepository,
	}
}

func (ges *groupExpenseServiceImpl) CreateDraft(ctx context.Context, request dto.NewGroupExpenseRequest) (dto.GroupExpenseResponse, error) {
	if err := ges.validate(ctx, request); err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	groupExpense := mapper.GroupExpenseRequestToEntity(request)

	insertedGroupExpense, err := ges.groupExpenseRepository.Insert(ctx, groupExpense)
	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	return mapper.GroupExpenseToResponse(insertedGroupExpense), nil
}

func (ges *groupExpenseServiceImpl) GetAllCreated(ctx context.Context, profileID uuid.UUID) ([]dto.GroupExpenseResponse, error) {
	spec := crud.Specification[entity.GroupExpense]{}
	spec.Model.CreatorProfileID = profileID
	spec.PreloadRelations = []string{"Items", "OtherFees"}

	groupExpenses, err := ges.groupExpenseRepository.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(groupExpenses, mapper.GroupExpenseToResponse), nil
}

func (ges *groupExpenseServiceImpl) GetDetails(ctx context.Context, id uuid.UUID) (dto.GroupExpenseResponse, error) {
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

	return mapper.GroupExpenseToResponse(groupExpense), nil
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
		// TODO: explore cleaner way
		groupExpense.Participants = nil // Prevent GORM updating child, already synced above

		updatedGroupExpense, err := ges.groupExpenseRepository.Update(ctx, groupExpense)
		if err != nil {
			return err
		}

		updatedGroupExpense.Participants = updatedGroupExpenseParticipants

		response = mapper.GroupExpenseToResponse(updatedGroupExpense)

		return nil
	})

	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	return response, nil
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

func (ges *groupExpenseServiceImpl) validate(ctx context.Context, request dto.NewGroupExpenseRequest) error {
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

	return nil
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

func (ges *groupExpenseServiceImpl) GetUnconfirmedGroupExpenseForUpdate(ctx context.Context, profileID, id uuid.UUID) (entity.GroupExpense, error) {
	spec := crud.Specification[entity.GroupExpense]{}
	spec.Model.ID = id
	spec.Model.CreatorProfileID = profileID
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
