package debt

import (
	"fmt"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/helper"
	"github.com/itsLeonB/ezutil"
)

type receivingAnonDebtCalculator struct {
	action appconstant.Action
}

func newReceivingAnonDebtCalculator() AnonymousDebtCalculator {
	return &receivingAnonDebtCalculator{
		action: appconstant.ReceiveAction,
	}
}

func (dc *receivingAnonDebtCalculator) GetAction() appconstant.Action {
	return dc.action
}

func (dc *receivingAnonDebtCalculator) MapRequestToEntity(request dto.NewDebtTransactionRequest) entity.DebtTransaction {
	return entity.DebtTransaction{
		LenderProfileID:   request.UserProfileID,
		BorrowerProfileID: request.FriendProfileID,
		Type:              appconstant.Repay,
		Amount:            request.Amount,
		TransferMethodID:  request.TransferMethodID,
		Description:       request.Description,
	}
}

func (dc *receivingAnonDebtCalculator) MapEntityToResponse(debtTransaction entity.DebtTransaction) dto.DebtTransactionResponse {
	return dto.DebtTransactionResponse{
		ID:             debtTransaction.ID,
		ProfileID:      debtTransaction.BorrowerProfileID,
		Type:           debtTransaction.Type,
		Amount:         debtTransaction.Amount,
		TransferMethod: debtTransaction.TransferMethod.Display,
		Description:    debtTransaction.Description,
		CreatedAt:      debtTransaction.CreatedAt,
		UpdatedAt:      debtTransaction.UpdatedAt,
		DeletedAt:      debtTransaction.DeletedAt.Time,
	}
}

func (dc *receivingAnonDebtCalculator) Validate(newTransaction entity.DebtTransaction, allTransactions []entity.DebtTransaction) error {
	userAmount, friendAmount := helper.GetAmountSumsFromDebtTransactions(
		newTransaction.LenderProfileID,
		newTransaction.BorrowerProfileID,
		allTransactions,
	)

	toReceiveLeftAmount := friendAmount.Sub(userAmount)

	if toReceiveLeftAmount.Compare(newTransaction.Amount) < 0 {
		return ezutil.ValidationError(fmt.Sprintf(
			"cannot receive debt, amount in user: %s, amount in friend: %s",
			userAmount,
			friendAmount,
		))
	}

	return nil
}
