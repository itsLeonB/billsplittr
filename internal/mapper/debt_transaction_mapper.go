package mapper

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil"
)

func NewDebtTransactionRequestToEntity(
	userProfileID uuid.UUID,
	request dto.NewDebtTransactionRequest,
) (entity.DebtTransaction, error) {
	lenderProfileID, borrowerProfileID, err := selectLenderBorrowerProfileID(userProfileID, request)
	if err != nil {
		return entity.DebtTransaction{}, err
	}

	return entity.DebtTransaction{
		LenderProfileID:   lenderProfileID,
		BorrowerProfileID: borrowerProfileID,
		Type:              appconstant.Lend,
		Amount:            request.Amount,
		TransferMethodID:  request.TransferMethodID,
		Description:       request.Description,
	}, nil
}

func selectLenderBorrowerProfileID(
	userProfileID uuid.UUID,
	request dto.NewDebtTransactionRequest,
) (uuid.UUID, uuid.UUID, error) {
	switch request.Action {
	case appconstant.BorrowAction:
		return request.FriendProfileID, userProfileID, nil
	case appconstant.LendAction:
		return userProfileID, request.FriendProfileID, nil
	default:
		return uuid.Nil, uuid.Nil, ezutil.UnprocessableEntityError(fmt.Sprintf("unsupported action: %s", request.Action))
	}
}
