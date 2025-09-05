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
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
)

type expenseItemServiceImpl struct {
	transactor             crud.Transactor
	groupExpenseRepository repository.GroupExpenseRepository
	expenseItemRepository  repository.ExpenseItemRepository
	groupExpenseSvc        GroupExpenseService
}

func NewExpenseItemService(
	transactor crud.Transactor,
	groupExpenseRepository repository.GroupExpenseRepository,
	expenseItemRepository repository.ExpenseItemRepository,
	groupExpenseSvc GroupExpenseService,
) ExpenseItemService {
	return &expenseItemServiceImpl{
		transactor,
		groupExpenseRepository,
		expenseItemRepository,
		groupExpenseSvc,
	}
}

func (ges *expenseItemServiceImpl) Add(ctx context.Context, request dto.NewExpenseItemRequest) (dto.ExpenseItemResponse, error) {
	var response dto.ExpenseItemResponse

	if !request.Amount.IsPositive() {
		return dto.ExpenseItemResponse{}, ungerr.UnprocessableEntityError(appconstant.ErrNonPositiveAmount)
	}

	err := ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		groupExpense, err := ges.groupExpenseSvc.GetUnconfirmedGroupExpenseForUpdate(ctx, request.ProfileID, request.GroupExpenseID)
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

		response = mapper.ExpenseItemToResponse(insertedItem)

		return nil
	})

	if err != nil {
		return dto.ExpenseItemResponse{}, err
	}

	return response, nil
}

func (ges *expenseItemServiceImpl) GetDetails(ctx context.Context, groupExpenseID, expenseItemID uuid.UUID) (dto.ExpenseItemResponse, error) {
	spec := crud.Specification[entity.ExpenseItem]{}
	spec.Model.ID = expenseItemID
	spec.Model.GroupExpenseID = groupExpenseID
	spec.PreloadRelations = []string{"Participants"}

	expenseItem, err := ges.getExpenseItemBySpec(ctx, spec)
	if err != nil {
		return dto.ExpenseItemResponse{}, err
	}

	return mapper.ExpenseItemToResponse(expenseItem), nil
}

func (ges *expenseItemServiceImpl) Update(ctx context.Context, request dto.UpdateExpenseItemRequest) (dto.ExpenseItemResponse, error) {
	var response dto.ExpenseItemResponse

	if !request.Amount.IsPositive() {
		return dto.ExpenseItemResponse{}, ungerr.UnprocessableEntityError(appconstant.ErrNonPositiveAmount)
	}

	err := ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		expenseItem, err := ges.getExpenseItemByIDForUpdate(ctx, request.ID, request.GroupExpenseID)
		if err != nil {
			return err
		}

		if ezutil.CompareUUID(request.GroupExpenseID, expenseItem.GroupExpenseID) != 0 {
			return ungerr.UnprocessableEntityError("mismatched group expense ID")
		}

		groupExpense, err := ges.groupExpenseSvc.GetUnconfirmedGroupExpenseForUpdate(ctx, request.ProfileID, expenseItem.GroupExpenseID)
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

		response = mapper.ExpenseItemToResponse(updatedExpenseItem)

		return nil
	})

	if err != nil {
		return dto.ExpenseItemResponse{}, err
	}

	return response, nil
}

func (ges *expenseItemServiceImpl) Remove(ctx context.Context, profileID, id, groupExpenseID uuid.UUID) error {
	return ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		groupExpense, err := ges.groupExpenseSvc.GetUnconfirmedGroupExpenseForUpdate(ctx, profileID, groupExpenseID)
		if err != nil {
			return err
		}

		expenseItem, err := ges.getExpenseItemByIDForUpdate(ctx, id, groupExpenseID)
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

func (ges *expenseItemServiceImpl) getExpenseItemByIDForUpdate(ctx context.Context, expenseItemID, groupExpenseID uuid.UUID) (entity.ExpenseItem, error) {
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

func (ges *expenseItemServiceImpl) getExpenseItemBySpec(ctx context.Context, spec crud.Specification[entity.ExpenseItem]) (entity.ExpenseItem, error) {
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
