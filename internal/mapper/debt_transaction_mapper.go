package mapper

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/shopspring/decimal"
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
		Action:         transaction.Action,
		Amount:         transaction.Amount,
		TransferMethod: transaction.TransferMethod.Display,
		Description:    transaction.Description,
		CreatedAt:      transaction.CreatedAt,
		UpdatedAt:      transaction.UpdatedAt,
		DeletedAt:      transaction.DeletedAt.Time,
	}
}

func GetDebtTransactionSimpleMapper(userProfileID uuid.UUID) func(transaction entity.DebtTransaction) dto.DebtTransactionResponse {
	return func(transaction entity.DebtTransaction) dto.DebtTransactionResponse {
		return DebtTransactionToResponse(userProfileID, transaction)
	}
}

func MapToFriendBalanceSummary(userProfileID uuid.UUID, debtTransactions []entity.DebtTransaction) dto.FriendBalance {
	totalOwedToYou, totalYouOwe := decimal.Zero, decimal.Zero

	for _, transaction := range debtTransactions {
		if transaction.LenderProfileID == userProfileID && transaction.Type == appconstant.Lend {
			totalOwedToYou = totalOwedToYou.Add(transaction.Amount)
		} else if transaction.LenderProfileID == userProfileID && transaction.Type == appconstant.Repay {
			totalOwedToYou = totalOwedToYou.Sub(transaction.Amount)
		} else if transaction.BorrowerProfileID == userProfileID && transaction.Type == appconstant.Lend {
			totalYouOwe = totalYouOwe.Add(transaction.Amount)
		} else if transaction.BorrowerProfileID == userProfileID && transaction.Type == appconstant.Repay {
			totalYouOwe = totalYouOwe.Sub(transaction.Amount)
		}
	}

	return dto.FriendBalance{
		TotalOwedToYou: totalOwedToYou,
		TotalYouOwe:    totalYouOwe,
		NetBalance:     totalOwedToYou.Sub(totalYouOwe),
		Currency:       appconstant.IDR,
	}
}
