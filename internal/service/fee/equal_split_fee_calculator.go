package fee

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

type equalSplitFeeCalculator struct {
	method appconstant.FeeCalculationMethod
}

func newEqualSplitFeeCalculator() FeeCalculator {
	return &equalSplitFeeCalculator{
		appconstant.EqualSplitFee,
	}
}

func (fc *equalSplitFeeCalculator) GetMethod() appconstant.FeeCalculationMethod {
	return fc.method
}

func (fc *equalSplitFeeCalculator) Validate(fee entity.OtherFee, groupExpense entity.GroupExpense) error {
	if fee.ID == uuid.Nil {
		return eris.New("fee ID cannot be nil")
	}

	if fee.GroupExpenseID == uuid.Nil {
		return eris.New("group expense ID cannot be nil")
	}

	if fee.Amount.IsZero() {
		return eris.New("amount cannot be zero")
	}

	if len(groupExpense.Participants) < 1 {
		return eris.New("must have participants")
	}

	return nil
}

func (fc *equalSplitFeeCalculator) Split(fee entity.OtherFee, groupExpense entity.GroupExpense) []entity.FeeParticipant {
	participantsCount := len(groupExpense.Participants)
	feeParticipants := make([]entity.FeeParticipant, participantsCount)
	amountPerParticipant := fee.Amount.Div(decimal.NewFromInt(int64(participantsCount)))

	for i, expenseParticipant := range groupExpense.Participants {
		feeParticipants[i] = entity.FeeParticipant{
			OtherFeeID:  fee.ID,
			ProfileID:   expenseParticipant.ParticipantProfileID,
			ShareAmount: amountPerParticipant,
		}
	}

	return feeParticipants
}

func (fc *equalSplitFeeCalculator) GetInfo() dto.FeeCalculationMethodInfo {
	return dto.FeeCalculationMethodInfo{
		Name:        fc.method,
		Display:     "Equal split",
		Description: "Equally split the fee to all participants",
	}
}
