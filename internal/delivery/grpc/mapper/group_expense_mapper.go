package mapper

import (
	"github.com/itsLeonB/audit/gen/go/audit/v1"
	"github.com/itsLeonB/billsplittr-protos/gen/go/groupexpense/v1"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/gerpc"
	"github.com/rotisserie/eris"
	"golang.org/x/text/currency"
)

func toExpenseParticipantResponseProto(expenseParticipant dto.ExpenseParticipantResponse) *groupexpense.ExpenseParticipantResponse {
	return &groupexpense.ExpenseParticipantResponse{
		ProfileId:   expenseParticipant.ProfileID.String(),
		ShareAmount: ezutil.DecimalToMoney(expenseParticipant.ShareAmount, currency.IDR.String()),
	}
}

func ToGroupExpenseResponseProto(groupExpense dto.GroupExpenseResponse) (*groupexpense.GroupExpenseResponse, error) {
	feeResponses, err := ezutil.MapSliceWithError(groupExpense.OtherFees, ToOtherFeeResponseProto)
	if err != nil {
		return nil, err
	}

	status, err := toExpenseStatusProto(groupExpense.Status)
	if err != nil {
		return nil, err
	}

	bill, err := ToExpenseBillResponseProto(groupExpense.Bill)
	if err != nil {
		return nil, err
	}

	return &groupexpense.GroupExpenseResponse{
		CreatorProfileId:        groupExpense.CreatorProfileID.String(),
		PayerProfileId:          groupExpense.PayerProfileID.String(),
		TotalAmount:             ezutil.DecimalToMoney(groupExpense.TotalAmount, currency.IDR.String()),
		ItemsTotal:              ezutil.DecimalToMoney(groupExpense.ItemsTotal, currency.IDR.String()),
		FeesTotal:               ezutil.DecimalToMoney(groupExpense.FeesTotal, currency.IDR.String()),
		Subtotal:                ezutil.DecimalToMoney(groupExpense.Subtotal, currency.IDR.String()),
		Description:             groupExpense.Description,
		IsConfirmed:             groupExpense.Confirmed,
		IsParticipantsConfirmed: groupExpense.ParticipantsConfirmed,
		Status:                  status,
		Items:                   ezutil.MapSlice(groupExpense.Items, ToExpenseItemResponseProto),
		OtherFees:               feeResponses,
		Participants:            ezutil.MapSlice(groupExpense.Participants, toExpenseParticipantResponseProto),
		ExpenseBill:             bill,
		AuditMetadata: &audit.Metadata{
			Id:        groupExpense.ID.String(),
			CreatedAt: gerpc.NullableTimeToProto(groupExpense.CreatedAt),
			UpdatedAt: gerpc.NullableTimeToProto(groupExpense.UpdatedAt),
			DeletedAt: gerpc.NullableTimeToProto(groupExpense.DeletedAt),
		},
	}, nil
}

func toExpenseStatusProto(status appconstant.ExpenseStatus) (groupexpense.GroupExpenseResponse_Status, error) {
	switch status {
	case appconstant.DraftExpense:
		return groupexpense.GroupExpenseResponse_STATUS_DRAFT, nil
	case appconstant.ReadyExpense:
		return groupexpense.GroupExpenseResponse_STATUS_READY, nil
	case appconstant.ConfirmedExpense:
		return groupexpense.GroupExpenseResponse_STATUS_CONFIRMED, nil
	default:
		return groupexpense.GroupExpenseResponse_STATUS_UNSPECIFIED, eris.Errorf("unspecified expense status constant: %s", status)
	}
}
