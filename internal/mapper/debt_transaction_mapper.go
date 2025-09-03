package mapper

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/drex-protos/gen/go/debt/v1"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"
)

func MapToFriendBalanceSummary(userProfileID uuid.UUID, debtTransactions []dto.DebtTransactionResponse) dto.FriendBalance {
	totalOwedToYou, totalYouOwe := decimal.Zero, decimal.Zero

	for _, transaction := range debtTransactions {
		switch transaction.Type {
		case appconstant.Lend:
			switch transaction.Action {
			case appconstant.LendAction: // You lent money
				totalOwedToYou = totalOwedToYou.Add(transaction.Amount)
			case appconstant.BorrowAction: // You borrowed money
				totalYouOwe = totalYouOwe.Add(transaction.Amount)
			}
		case appconstant.Repay:
			switch transaction.Action {
			case appconstant.ReceiveAction: // You received repayment
				totalOwedToYou = totalOwedToYou.Sub(transaction.Amount)
			case appconstant.ReturnAction: // You returned money
				totalYouOwe = totalYouOwe.Sub(transaction.Amount)
			}
		}
	}

	return dto.FriendBalance{
		TotalOwedToYou: totalOwedToYou,
		TotalYouOwe:    totalYouOwe,
		NetBalance:     totalOwedToYou.Sub(totalYouOwe),
		CurrencyCode:   currency.IDR.String(),
	}
}

func ToProtoTransactionAction(ta appconstant.Action) (debt.TransactionAction, error) {
	switch ta {
	case appconstant.BorrowAction:
		return debt.TransactionAction_TRANSACTION_ACTION_BORROW, nil
	case appconstant.LendAction:
		return debt.TransactionAction_TRANSACTION_ACTION_LEND, nil
	case appconstant.ReceiveAction:
		return debt.TransactionAction_TRANSACTION_ACTION_RECEIVE, nil
	case appconstant.ReturnAction:
		return debt.TransactionAction_TRANSACTION_ACTION_RETURN, nil
	default:
		return debt.TransactionAction_TRANSACTION_ACTION_UNSPECIFIED, eris.Errorf("undefined TransactionAction constant: %s", ta)
	}
}

func FromProtoTransactionAction(ta debt.TransactionAction) (appconstant.Action, error) {
	switch ta {
	case debt.TransactionAction_TRANSACTION_ACTION_BORROW:
		return appconstant.BorrowAction, nil
	case debt.TransactionAction_TRANSACTION_ACTION_LEND:
		return appconstant.LendAction, nil
	case debt.TransactionAction_TRANSACTION_ACTION_RECEIVE:
		return appconstant.ReceiveAction, nil
	case debt.TransactionAction_TRANSACTION_ACTION_RETURN:
		return appconstant.ReturnAction, nil
	default:
		return "", eris.Errorf("undefined TransactionAction enum: %s", ta)
	}
}

func FromProtoTransactionType(ta debt.TransactionType) (appconstant.DebtTransactionType, error) {
	switch ta {
	case debt.TransactionType_TRANSACTION_TYPE_LEND:
		return appconstant.Lend, nil
	case debt.TransactionType_TRANSACTION_TYPE_REPAY:
		return appconstant.Repay, nil
	default:
		return "", eris.Errorf("undefined TransactionType enum: %s", ta)
	}
}

func FromProtoDebtTransactionResponse(trx *debt.TransactionResponse) (dto.DebtTransactionResponse, error) {
	if trx == nil {
		return dto.DebtTransactionResponse{}, eris.New("transaction from response is nil")
	}

	id, err := ezutil.Parse[uuid.UUID](trx.GetId())
	if err != nil {
		return dto.DebtTransactionResponse{}, err
	}

	profileID, err := ezutil.Parse[uuid.UUID](trx.GetProfileId())
	if err != nil {
		return dto.DebtTransactionResponse{}, err
	}

	trxType, err := FromProtoTransactionType(trx.GetType())
	if err != nil {
		return dto.DebtTransactionResponse{}, err
	}

	trxAction, err := FromProtoTransactionAction(trx.GetAction())
	if err != nil {
		return dto.DebtTransactionResponse{}, err
	}

	return dto.DebtTransactionResponse{
		ID:             id,
		ProfileID:      profileID,
		Type:           trxType,
		Action:         trxAction,
		Amount:         ezutil.MoneyToDecimal(trx.GetAmount()),
		TransferMethod: trx.GetTransferMethod(),
		Description:    trx.GetDescription(),
		CreatedAt:      ezutil.FromProtoTime(trx.GetCreatedAt()),
		UpdatedAt:      ezutil.FromProtoTime(trx.GetUpdatedAt()),
		DeletedAt:      ezutil.FromProtoTime(trx.GetDeletedAt()),
	}, nil
}
