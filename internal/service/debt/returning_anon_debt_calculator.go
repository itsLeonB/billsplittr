package debt

import (
	"fmt"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/helper"
	"github.com/itsLeonB/ungerr"
)

type returningAnonDebtCalculator struct {
	action appconstant.Action
}

func newReturningAnonDebtCalculator() AnonymousDebtCalculator {
	return &returningAnonDebtCalculator{
		action: appconstant.ReturnAction,
	}
}

func (dc *returningAnonDebtCalculator) GetAction() appconstant.Action {
	return dc.action
}

func (dc *returningAnonDebtCalculator) MapRequestToEntity(request dto.NewDebtTransactionRequest) entity.DebtTransaction {
	return entity.DebtTransaction{
		LenderProfileID:   request.FriendProfileID,
		BorrowerProfileID: request.UserProfileID,
		Action:            dc.action,
		Type:              appconstant.Repay,
		Amount:            request.Amount,
		TransferMethodID:  request.TransferMethodID,
		Description:       request.Description,
	}
}

func (dc *returningAnonDebtCalculator) MapEntityToResponse(debtTransaction entity.DebtTransaction) dto.DebtTransactionResponse {
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

func (dc *returningAnonDebtCalculator) Validate(newTransaction entity.DebtTransaction, allTransactions []entity.DebtTransaction) error {
	userAmount, friendAmount := helper.GetAmountSumsFromDebtTransactions(
		newTransaction.BorrowerProfileID,
		newTransaction.LenderProfileID,
		allTransactions,
	)

	toReturnLeftAmount := userAmount.Sub(friendAmount)

	if toReturnLeftAmount.Compare(newTransaction.Amount) < 0 {
		return ungerr.ValidationError(fmt.Sprintf(
			"cannot return debt, amount in user: %s, amount in friend: %s",
			userAmount,
			friendAmount,
		))
	}

	return nil
}
