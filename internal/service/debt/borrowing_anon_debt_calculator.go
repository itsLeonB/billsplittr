package debt

import (
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
)

type borrowingAnonDebtCalculator struct {
	action appconstant.Action
}

func newBorrowingAnonDebtCalculator() AnonymousDebtCalculator {
	return &borrowingAnonDebtCalculator{
		action: appconstant.BorrowAction,
	}
}

func (dc *borrowingAnonDebtCalculator) GetAction() appconstant.Action {
	return dc.action
}

func (dc *borrowingAnonDebtCalculator) MapRequestToEntity(request dto.NewDebtTransactionRequest) entity.DebtTransaction {
	return entity.DebtTransaction{
		LenderProfileID:   request.FriendProfileID,
		BorrowerProfileID: request.UserProfileID,
		Type:              appconstant.Lend,
		Action:            dc.action,
		Amount:            request.Amount,
		TransferMethodID:  request.TransferMethodID,
		Description:       request.Description,
	}
}

func (dc *borrowingAnonDebtCalculator) MapEntityToResponse(debtTransaction entity.DebtTransaction) dto.DebtTransactionResponse {
	return dto.DebtTransactionResponse{
		ID:             debtTransaction.ID,
		ProfileID:      debtTransaction.LenderProfileID,
		Type:           debtTransaction.Type,
		Action:         debtTransaction.Action,
		Amount:         debtTransaction.Amount,
		TransferMethod: debtTransaction.TransferMethod.Display,
		Description:    debtTransaction.Description,
		CreatedAt:      debtTransaction.CreatedAt,
		UpdatedAt:      debtTransaction.UpdatedAt,
		DeletedAt:      debtTransaction.DeletedAt.Time,
	}
}

func (dc *borrowingAnonDebtCalculator) Validate(newTransaction entity.DebtTransaction, allTransactions []entity.DebtTransaction) error {
	// Currently does not validate stuff
	// User can record borrow of any amount for anonymous friend
	return nil
}
