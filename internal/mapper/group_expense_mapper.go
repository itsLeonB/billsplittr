package mapper

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil"
)

func GroupExpenseRequestToEntity(request dto.NewGroupExpenseRequest) entity.GroupExpense {
	return entity.GroupExpense{
		PayerProfileID:   request.PayerProfileID,
		TotalAmount:      request.TotalAmount,
		Description:      request.Description,
		Items:            ezutil.MapSlice(request.Items, expenseItemRequestToEntity),
		OtherFees:        ezutil.MapSlice(request.OtherFees, otherFeeRequestToEntity),
		CreatorProfileID: request.CreatedByProfileID,
	}
}

func GetGroupExpenseSimpleMapper(userProfileID uuid.UUID) func(groupExpense entity.GroupExpense) dto.GroupExpenseResponse {
	return func(groupExpense entity.GroupExpense) dto.GroupExpenseResponse {
		return GroupExpenseToResponse(groupExpense, userProfileID)
	}
}

func GroupExpenseToResponse(groupExpense entity.GroupExpense, userProfileID uuid.UUID) dto.GroupExpenseResponse {
	return dto.GroupExpenseResponse{
		ID:                    groupExpense.ID,
		PayerProfileID:        groupExpense.PayerProfileID,
		PayerName:             groupExpense.PayerProfile.Name,
		PaidByUser:            groupExpense.PayerProfileID == userProfileID,
		TotalAmount:           groupExpense.TotalAmount,
		Description:           groupExpense.Description,
		Items:                 ezutil.MapSlice(groupExpense.Items, getExpenseItemSimpleMapper(userProfileID)),
		OtherFees:             ezutil.MapSlice(groupExpense.OtherFees, otherFeeToResponse),
		CreatorProfileID:      groupExpense.CreatorProfileID,
		CreatorName:           groupExpense.CreatorProfile.Name,
		CreatedByUser:         groupExpense.CreatorProfileID == userProfileID,
		Confirmed:             groupExpense.Confirmed,
		ParticipantsConfirmed: groupExpense.ParticipantsConfirmed,
		CreatedAt:             groupExpense.CreatedAt,
		UpdatedAt:             groupExpense.UpdatedAt,
		DeletedAt:             groupExpense.DeletedAt.Time,
		Participants:          ezutil.MapSlice(groupExpense.Participants, getExpenseParticipantSimpleMapper(userProfileID)),
	}
}

func getExpenseItemSimpleMapper(userProfileID uuid.UUID) func(item entity.ExpenseItem) dto.ExpenseItemResponse {
	return func(item entity.ExpenseItem) dto.ExpenseItemResponse {
		return ExpenseItemToResponse(item, userProfileID)
	}
}

func ExpenseItemToResponse(item entity.ExpenseItem, userProfileID uuid.UUID) dto.ExpenseItemResponse {
	return dto.ExpenseItemResponse{
		ID:             item.ID,
		GroupExpenseID: item.GroupExpenseID,
		Name:           item.Name,
		Amount:         item.Amount,
		Quantity:       item.Quantity,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		DeletedAt:      item.DeletedAt.Time,
		Participants:   ezutil.MapSlice(item.Participants, getItemParticipantSimpleMapper(userProfileID)),
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

func getItemParticipantSimpleMapper(userProfileID uuid.UUID) func(itemParticipant entity.ItemParticipant) dto.ItemParticipantResponse {
	return func(itemParticipant entity.ItemParticipant) dto.ItemParticipantResponse {
		return itemParticipantToResponse(itemParticipant, userProfileID)
	}
}

func itemParticipantToResponse(itemParticipant entity.ItemParticipant, userProfileID uuid.UUID) dto.ItemParticipantResponse {
	return dto.ItemParticipantResponse{
		ProfileName: itemParticipant.Profile.Name,
		ProfileID:   itemParticipant.ProfileID,
		Share:       itemParticipant.Share,
		IsUser:      itemParticipant.ProfileID == userProfileID,
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

func PatchExpenseItemWithRequest(expenseItem entity.ExpenseItem, request dto.UpdateExpenseItemRequest) entity.ExpenseItem {
	expenseItem.Name = request.Name
	expenseItem.Amount = request.Amount
	expenseItem.Quantity = request.Quantity
	return expenseItem
}

func ItemParticipantRequestToEntity(itemParticipant dto.NewItemParticipantRequest) entity.ItemParticipant {
	return entity.ItemParticipant{
		ProfileID: itemParticipant.ProfileID,
		Share:     itemParticipant.Share,
	}
}

func ExpenseParticipantToResponse(expenseParticipant entity.ExpenseParticipant, userProfileID uuid.UUID) dto.ExpenseParticipantResponse {
	return dto.ExpenseParticipantResponse{
		ProfileName: expenseParticipant.Profile.Name,
		ProfileID:   expenseParticipant.ParticipantProfileID,
		ShareAmount: expenseParticipant.ShareAmount,
		IsUser:      expenseParticipant.ParticipantProfileID == userProfileID,
	}
}

func getExpenseParticipantSimpleMapper(userProfileID uuid.UUID) func(entity.ExpenseParticipant) dto.ExpenseParticipantResponse {
	return func(ep entity.ExpenseParticipant) dto.ExpenseParticipantResponse {
		return ExpenseParticipantToResponse(ep, userProfileID)
	}
}

func GroupExpenseToDebtTransactions(groupExpense entity.GroupExpense, transferMethodID uuid.UUID) []entity.DebtTransaction {
	action := appconstant.BorrowAction
	if groupExpense.PayerProfileID == groupExpense.CreatorProfileID {
		action = appconstant.LendAction
	}

	debtTransactions := make([]entity.DebtTransaction, 0, len(groupExpense.Participants))
	for _, participant := range groupExpense.Participants {
		if groupExpense.PayerProfileID == participant.ParticipantProfileID {
			continue
		}
		debtTransactions = append(debtTransactions, entity.DebtTransaction{
			LenderProfileID:   groupExpense.PayerProfileID,
			BorrowerProfileID: participant.ParticipantProfileID,
			Type:              appconstant.Lend,
			Action:            action,
			Amount:            participant.ShareAmount,
			TransferMethodID:  transferMethodID,
			Description:       fmt.Sprintf("Share for group expense: %s", groupExpense.Description),
		})
	}

	return debtTransactions
}
