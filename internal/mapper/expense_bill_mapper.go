package mapper

import (
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
)

func ExpenseBillToResponse(eb entity.ExpenseBill) dto.ExpenseBillResponse {
	return dto.ExpenseBillResponse{
		ID:               eb.ID,
		CreatorProfileID: eb.CreatorProfileID,
		PayerProfileID:   eb.PayerProfileID,
		GroupExpenseID:   eb.GroupExpenseID.UUID,
		Filename:         eb.ImageName,
		CreatedAt:        eb.CreatedAt,
		UpdatedAt:        eb.UpdatedAt,
		DeletedAt:        eb.DeletedAt.Time,
	}
}
