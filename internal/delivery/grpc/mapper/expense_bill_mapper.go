package mapper

import (
	"github.com/itsLeonB/audit/gen/go/audit/v1"
	"github.com/itsLeonB/billsplittr-protos/gen/go/expensebill/v1"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/gerpc"
)

func ToExpenseBillResponseProto(eb dto.ExpenseBillResponse) *expensebill.ExpenseBillResponse {
	return &expensebill.ExpenseBillResponse{
		ExpenseBill: &expensebill.ExpenseBill{
			CreatorProfileId: eb.CreatorProfileID.String(),
			PayerProfileId:   eb.PayerProfileID.String(),
			GroupExpenseId:   eb.GroupExpenseID.String(),
			ObjectKey:        eb.Filename,
		},
		AuditMetadata: &audit.Metadata{
			Id:        eb.ID.String(),
			CreatedAt: gerpc.NullableTimeToProto(eb.CreatedAt),
			UpdatedAt: gerpc.NullableTimeToProto(eb.UpdatedAt),
			DeletedAt: gerpc.NullableTimeToProto(eb.DeletedAt),
		},
	}
}
