package mapper

import (
	"github.com/itsLeonB/audit/gen/go/audit/v1"
	"github.com/itsLeonB/billsplittr-protos/gen/go/expensebill/v1"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/gerpc"
	"github.com/rotisserie/eris"
)

func ToExpenseBillResponseProto(eb dto.ExpenseBillResponse) (*expensebill.ExpenseBillResponse, error) {
	status, err := toBillStatusProto(eb.Status)
	if err != nil {
		return nil, err
	}

	return &expensebill.ExpenseBillResponse{
		ExpenseBill: &expensebill.ExpenseBill{
			CreatorProfileId: eb.CreatorProfileID.String(),
			PayerProfileId:   eb.PayerProfileID.String(),
			GroupExpenseId:   eb.GroupExpenseID.String(),
			ObjectKey:        eb.Filename,
			Status:           status,
		},
		AuditMetadata: &audit.Metadata{
			Id:        eb.ID.String(),
			CreatedAt: gerpc.NullableTimeToProto(eb.CreatedAt),
			UpdatedAt: gerpc.NullableTimeToProto(eb.UpdatedAt),
			DeletedAt: gerpc.NullableTimeToProto(eb.DeletedAt),
		},
	}, nil
}

func toBillStatusProto(status appconstant.BillStatus) (expensebill.ExpenseBill_Status, error) {
	switch status {
	case appconstant.PendingBill:
		return expensebill.ExpenseBill_STATUS_PENDING, nil
	case appconstant.ParsedBill:
		return expensebill.ExpenseBill_STATUS_PARSED, nil
	case appconstant.FailedBill:
		return expensebill.ExpenseBill_STATUS_FAILED, nil
	default:
		return expensebill.ExpenseBill_STATUS_UNSPECIFIED, eris.Errorf("unspecified bill status constant: %s", status)
	}
}
