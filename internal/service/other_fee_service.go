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
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/shopspring/decimal"
)

type otherFeeServiceImpl struct {
	transactor             crud.Transactor
	groupExpenseRepository repository.GroupExpenseRepository
	feeCalculatorRegistry  map[appconstant.FeeCalculationMethod]fee.FeeCalculator
	otherFeeRepository     repository.OtherFeeRepository
	groupExpenseSvc        GroupExpenseService
}

func NewOtherFeeService(
	transactor crud.Transactor,
	groupExpenseRepository repository.GroupExpenseRepository,
	otherFeeRepository repository.OtherFeeRepository,
	groupExpenseSvc GroupExpenseService,
) OtherFeeService {
	return &otherFeeServiceImpl{
		transactor,
		groupExpenseRepository,
		fee.NewFeeCalculatorRegistry(),
		otherFeeRepository,
		groupExpenseSvc,
	}
}

func (ges *otherFeeServiceImpl) Add(ctx context.Context, request dto.NewOtherFeeRequest) (dto.OtherFeeResponse, error) {
	var response dto.OtherFeeResponse

	if !request.Amount.IsPositive() {
		return dto.OtherFeeResponse{}, ungerr.UnprocessableEntityError(appconstant.ErrNonPositiveAmount)
	}

	err := ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		groupExpense, err := ges.groupExpenseSvc.GetUnconfirmedGroupExpenseForUpdate(ctx, request.ProfileID, request.GroupExpenseID)
		if err != nil {
			return err
		}

		fee := mapper.OtherFeeRequestToEntity(request.OtherFeeData)
		fee.GroupExpenseID = request.GroupExpenseID

		groupExpense.TotalAmount = groupExpense.TotalAmount.Add(fee.Amount)
		if _, err = ges.groupExpenseRepository.Update(ctx, groupExpense); err != nil {
			return err
		}

		insertedFee, err := ges.otherFeeRepository.Insert(ctx, fee)
		if err != nil {
			return err
		}

		response = mapper.OtherFeeToResponse(insertedFee)

		return nil
	})

	if err != nil {
		return dto.OtherFeeResponse{}, err
	}

	return response, nil
}

func (ges *otherFeeServiceImpl) Update(ctx context.Context, request dto.UpdateOtherFeeRequest) (dto.OtherFeeResponse, error) {
	var response dto.OtherFeeResponse

	if request.Amount.Cmp(decimal.Zero) <= 0 {
		return dto.OtherFeeResponse{}, ungerr.UnprocessableEntityError("amount must be more than 0")
	}

	err := ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		groupExpense, err := ges.groupExpenseSvc.GetUnconfirmedGroupExpenseForUpdate(ctx, request.ProfileID, request.GroupExpenseID)
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

		response = mapper.OtherFeeToResponse(updatedFee)

		return nil
	})

	if err != nil {
		return dto.OtherFeeResponse{}, err
	}

	return response, nil
}

func (ges *otherFeeServiceImpl) Remove(ctx context.Context, id, profileID, groupExpenseID uuid.UUID) error {
	return ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		groupExpense, err := ges.groupExpenseSvc.GetUnconfirmedGroupExpenseForUpdate(ctx, profileID, groupExpenseID)
		if err != nil {
			return err
		}

		spec := crud.Specification[entity.OtherFee]{}
		spec.Model.ID = id
		spec.Model.GroupExpenseID = groupExpenseID
		spec.ForUpdate = true
		otherFee, err := ges.otherFeeRepository.FindFirst(ctx, spec)
		if err != nil {
			return err
		}
		if otherFee.IsZero() {
			return ungerr.NotFoundError(fmt.Sprintf("other fee with ID: %s is not found", id))
		}
		if otherFee.IsDeleted() {
			return ungerr.UnprocessableEntityError(fmt.Sprintf("other fee with ID: %s is deleted", id))
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

func (ges *otherFeeServiceImpl) GetCalculationMethods() []dto.FeeCalculationMethodInfo {
	feeCalculationMethodInfos := make([]dto.FeeCalculationMethodInfo, 0, len(ges.feeCalculatorRegistry))
	for _, feeCalculator := range ges.feeCalculatorRegistry {
		feeCalculationMethodInfos = append(feeCalculationMethodInfos, feeCalculator.GetInfo())
	}

	return feeCalculationMethodInfos
}
