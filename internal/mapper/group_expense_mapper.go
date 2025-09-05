package mapper

import (
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil/v2"
)

func GroupExpenseRequestToEntity(request dto.NewGroupExpenseRequest) entity.GroupExpense {
	return entity.GroupExpense{
		PayerProfileID:   request.PayerProfileID,
		TotalAmount:      request.TotalAmount,
		Subtotal:         request.Subtotal,
		Description:      request.Description,
		Items:            ezutil.MapSlice(request.Items, expenseItemDataToEntity),
		OtherFees:        ezutil.MapSlice(request.OtherFees, OtherFeeRequestToEntity),
		CreatorProfileID: request.CreatorProfileID,
	}
}

func GroupExpenseToResponse(groupExpense entity.GroupExpense) dto.GroupExpenseResponse {
	return dto.GroupExpenseResponse{
		ID:                    groupExpense.ID,
		PayerProfileID:        groupExpense.PayerProfileID,
		TotalAmount:           groupExpense.TotalAmount,
		Subtotal:              groupExpense.Subtotal,
		Description:           groupExpense.Description,
		Items:                 ezutil.MapSlice(groupExpense.Items, ExpenseItemToResponse),
		OtherFees:             ezutil.MapSlice(groupExpense.OtherFees, OtherFeeToResponse),
		CreatorProfileID:      groupExpense.CreatorProfileID,
		Confirmed:             groupExpense.Confirmed,
		ParticipantsConfirmed: groupExpense.ParticipantsConfirmed,
		CreatedAt:             groupExpense.CreatedAt,
		UpdatedAt:             groupExpense.UpdatedAt,
		DeletedAt:             groupExpense.DeletedAt.Time,
		Participants:          ezutil.MapSlice(groupExpense.Participants, expenseParticipantToResponse),
	}
}

func expenseParticipantToResponse(expenseParticipant entity.ExpenseParticipant) dto.ExpenseParticipantResponse {
	return dto.ExpenseParticipantResponse{
		ProfileID:   expenseParticipant.ParticipantProfileID,
		ShareAmount: expenseParticipant.ShareAmount,
	}
}
