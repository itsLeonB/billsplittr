package mapper

import (
	"github.com/itsLeonB/billsplittr-protos/gen/go/domain/v1"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/gerpc"
	"golang.org/x/text/currency"
)

func toItemParticipantProto(itemParticipant dto.ItemParticipantData) *domain.ItemParticipant {
	return &domain.ItemParticipant{
		ProfileId: itemParticipant.ProfileID.String(),
		Share:     float32(itemParticipant.Share.InexactFloat64()),
	}
}

func ToExpenseItemResponseProto(item dto.ExpenseItemResponse) *domain.ExpenseItemResponse {
	return &domain.ExpenseItemResponse{
		GroupExpenseId: item.GroupExpenseID.String(),
		ExpenseItem: &domain.ExpenseItem{
			Name:         item.Name,
			Amount:       ezutil.DecimalToMoney(item.Amount, currency.IDR.String()),
			Quantity:     int64(item.Quantity),
			Participants: ezutil.MapSlice(item.Participants, toItemParticipantProto),
		},
		AuditMetadata: &domain.AuditMetadata{
			Id:        item.ID.String(),
			CreatedAt: gerpc.NullableTimeToProto(item.CreatedAt),
			UpdatedAt: gerpc.NullableTimeToProto(item.UpdatedAt),
			DeletedAt: gerpc.NullableTimeToProto(item.DeletedAt),
		},
	}
}

func FromExpenseItemProto(item *domain.ExpenseItem) dto.ExpenseItemData {
	return dto.ExpenseItemData{
		Name:     item.GetName(),
		Amount:   ezutil.MoneyToDecimal(item.GetAmount()),
		Quantity: int(item.GetQuantity()),
	}
}
