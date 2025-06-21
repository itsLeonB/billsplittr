package mapper

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
)

func DebtTransactionToResponse(userProfileID uuid.UUID, transaction entity.DebtTransaction) dto.DebtTransactionResponse {
	var profileID uuid.UUID
	if userProfileID == transaction.BorrowerProfileID && userProfileID != transaction.LenderProfileID {
		profileID = transaction.LenderProfileID
	} else if userProfileID == transaction.LenderProfileID && userProfileID != transaction.BorrowerProfileID {
		profileID = transaction.BorrowerProfileID
	}

	return dto.DebtTransactionResponse{
		ID:             transaction.ID,
		ProfileID:      profileID,
		Type:           transaction.Type,
		Amount:         transaction.Amount,
		TransferMethod: transaction.TransferMethod.Display,
		Description:    transaction.Description,
		CreatedAt:      transaction.CreatedAt,
		UpdatedAt:      transaction.UpdatedAt,
		DeletedAt:      transaction.DeletedAt.Time,
	}
}
