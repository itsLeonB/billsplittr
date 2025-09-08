package mapper

import (
	"github.com/itsLeonB/audit/gen/go/audit/v1"
	"github.com/itsLeonB/billsplittr-protos/gen/go/expenseitem/v1"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/gerpc"
	"golang.org/x/text/currency"
)

func toItemParticipantProto(itemParticipant dto.ItemParticipantData) *expenseitem.ItemParticipant {
	return &expenseitem.ItemParticipant{
		ProfileId: itemParticipant.ProfileID.String(),
		Share:     itemParticipant.Share.InexactFloat64(),
	}
}

func ToExpenseItemResponseProto(item dto.ExpenseItemResponse) *expenseitem.ExpenseItemResponse {
	return &expenseitem.ExpenseItemResponse{
		GroupExpenseId: item.GroupExpenseID.String(),
		ExpenseItem: &expenseitem.ExpenseItem{
			Name:         item.Name,
			Amount:       ezutil.DecimalToMoney(item.Amount, currency.IDR.String()),
			Quantity:     int64(item.Quantity),
			Participants: ezutil.MapSlice(item.Participants, toItemParticipantProto),
		},
		AuditMetadata: &audit.Metadata{
			Id:        item.ID.String(),
			CreatedAt: gerpc.NullableTimeToProto(item.CreatedAt),
			UpdatedAt: gerpc.NullableTimeToProto(item.UpdatedAt),
			DeletedAt: gerpc.NullableTimeToProto(item.DeletedAt),
		},
	}
}

func FromExpenseItemProto(item *expenseitem.ExpenseItem) dto.ExpenseItemData {
	return dto.ExpenseItemData{
		Name:     item.GetName(),
		Amount:   ezutil.MoneyToDecimal(item.GetAmount()),
		Quantity: int(item.GetQuantity()),
	}
}
