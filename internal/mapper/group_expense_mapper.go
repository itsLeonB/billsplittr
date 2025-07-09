package mapper

import (
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil"
)

func GroupExpenseRequestToEntity(request dto.NewGroupExpenseRequest) entity.GroupExpense {
	return entity.GroupExpense{
		PayerProfileID:     request.PayerProfileID,
		TotalAmount:        request.TotalAmount,
		Description:        request.Description,
		Items:              ezutil.MapSlice(request.Items, expenseItemRequestToEntity),
		OtherFees:          ezutil.MapSlice(request.OtherFees, otherFeeRequestToEntity),
		CreatedByProfileID: request.CreatedByProfileID,
	}
}

func GroupExpenseToResponse(groupExpense entity.GroupExpense) dto.GroupExpenseResponse {
	return dto.GroupExpenseResponse{
		ID:                 groupExpense.ID,
		PayerProfileID:     groupExpense.PayerProfileID,
		TotalAmount:        groupExpense.TotalAmount,
		Description:        groupExpense.Description,
		Items:              ezutil.MapSlice(groupExpense.Items, expenseItemToResponse),
		OtherFees:          ezutil.MapSlice(groupExpense.OtherFees, otherFeeToResponse),
		CreatedByProfileID: groupExpense.CreatedByProfileID,
		CreatedAt:          groupExpense.CreatedAt,
		UpdatedAt:          groupExpense.UpdatedAt,
		DeletedAt:          groupExpense.DeletedAt.Time,
	}
}

func expenseItemToResponse(item entity.ExpenseItem) dto.ExpenseItemResponse {
	return dto.ExpenseItemResponse{
		ID:        item.ID,
		Name:      item.Name,
		Amount:    item.Amount,
		Quantity:  item.Quantity,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
		DeletedAt: item.DeletedAt.Time,
	}
}

func otherFeeToResponse(fee entity.OtherFee) dto.OtherFeeResponse {
	return dto.OtherFeeResponse{
		ID:        fee.ID,
		Name:      fee.Name,
		Amount:    fee.Amount,
		CreatedAt: fee.CreatedAt,
		UpdatedAt: fee.UpdatedAt,
		DeletedAt: fee.DeletedAt.Time,
	}
}

func expenseItemRequestToEntity(request dto.NewExpenseItemRequest) entity.ExpenseItem {
	return entity.ExpenseItem{
		Name:     request.Name,
		Amount:   request.Amount,
		Quantity: request.Quantity,
	}
}

func otherFeeRequestToEntity(request dto.NewOtherFeeRequest) entity.OtherFee {
	return entity.OtherFee{
		Name:   request.Name,
		Amount: request.Amount,
	}
}
