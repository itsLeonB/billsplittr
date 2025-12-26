package mapper

import (
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil/v2"
)

func GroupExpenseRequestToEntity(request dto.NewGroupExpenseRequest) entity.GroupExpense {
	status := appconstant.DraftExpense
	if len(request.Items) > 0 {
		status = appconstant.ReadyExpense
	}
	return entity.GroupExpense{
		PayerProfileID:   request.PayerProfileID,
		TotalAmount:      request.TotalAmount,
		Subtotal:         request.Subtotal,
		ItemsTotal:       request.Subtotal,
		FeesTotal:        request.TotalAmount.Sub(request.Subtotal),
		Description:      request.Description,
		Items:            ezutil.MapSlice(request.Items, expenseItemDataToEntity),
		OtherFees:        ezutil.MapSlice(request.OtherFees, OtherFeeRequestToEntity),
		CreatorProfileID: request.CreatorProfileID,
		Status:           status,
	}
}

func GroupExpenseToResponse(groupExpense entity.GroupExpense) dto.GroupExpenseResponse {
	return dto.GroupExpenseResponse{
		ID:               groupExpense.ID,
		PayerProfileID:   groupExpense.PayerProfileID,
		TotalAmount:      groupExpense.TotalAmount,
		Subtotal:         groupExpense.Subtotal,
		ItemsTotal:       groupExpense.ItemsTotal,
		FeesTotal:        groupExpense.FeesTotal,
		Description:      groupExpense.Description,
		CreatorProfileID: groupExpense.CreatorProfileID,
		Confirmed:        groupExpense.Confirmed,
		Status:           groupExpense.Status,
		CreatedAt:        groupExpense.CreatedAt,
		UpdatedAt:        groupExpense.UpdatedAt,
		DeletedAt:        groupExpense.DeletedAt.Time,
		Items:            ezutil.MapSlice(groupExpense.Items, ExpenseItemToResponse),
		OtherFees:        ezutil.MapSlice(groupExpense.OtherFees, OtherFeeToResponse),
		Participants:     ezutil.MapSlice(groupExpense.Participants, expenseParticipantToResponse),
		Bill:             ExpenseBillToResponse(groupExpense.Bill),
	}
}

func expenseParticipantToResponse(expenseParticipant entity.ExpenseParticipant) dto.ExpenseParticipantResponse {
	return dto.ExpenseParticipantResponse{
		ProfileID:   expenseParticipant.ParticipantProfileID,
		ShareAmount: expenseParticipant.ShareAmount,
	}
}
