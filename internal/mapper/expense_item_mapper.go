package mapper

import (
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil/v2"
)

func ExpenseItemToResponse(item entity.ExpenseItem) dto.ExpenseItemResponse {
	return dto.ExpenseItemResponse{
		ID:             item.ID,
		GroupExpenseID: item.GroupExpenseID,
		Name:           item.Name,
		Amount:         item.Amount,
		Quantity:       item.Quantity,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		DeletedAt:      item.DeletedAt.Time,
		Participants:   ezutil.MapSlice(item.Participants, itemParticipantToResponse),
	}
}

func itemParticipantToResponse(itemParticipant entity.ItemParticipant) dto.ItemParticipantData {
	return dto.ItemParticipantData{
		ProfileID: itemParticipant.ProfileID,
		Share:     itemParticipant.Share,
	}
}

func expenseItemDataToEntity(data dto.ExpenseItemData) entity.ExpenseItem {
	return entity.ExpenseItem{
		Name:     data.Name,
		Amount:   data.Amount,
		Quantity: data.Quantity,
	}
}

func ExpenseItemRequestToEntity(request dto.NewExpenseItemRequest) entity.ExpenseItem {
	return entity.ExpenseItem{
		GroupExpenseID: request.GroupExpenseID,
		Name:           request.Name,
		Amount:         request.Amount,
		Quantity:       request.Quantity,
	}
}

func PatchExpenseItemWithRequest(expenseItem entity.ExpenseItem, request dto.UpdateExpenseItemRequest) entity.ExpenseItem {
	expenseItem.Name = request.Name
	expenseItem.Amount = request.Amount
	expenseItem.Quantity = request.Quantity
	return expenseItem
}

func ItemParticipantRequestToEntity(itemParticipant dto.ItemParticipantData) entity.ItemParticipant {
	return entity.ItemParticipant{
		ProfileID: itemParticipant.ProfileID,
		Share:     itemParticipant.Share,
	}
}
