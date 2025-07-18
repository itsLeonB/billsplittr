package fee

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/rotisserie/eris"
)

type itemizedSplitFeeCalculator struct {
	method appconstant.FeeCalculationMethod
}

func newItemizedSplitFeeCalculator() FeeCalculator {
	return &itemizedSplitFeeCalculator{
		appconstant.ItemizedSplitFee,
	}
}

func (fc *itemizedSplitFeeCalculator) GetMethod() appconstant.FeeCalculationMethod {
	return fc.method
}

func (fc *itemizedSplitFeeCalculator) Validate(fee entity.OtherFee, groupExpense entity.GroupExpense) error {
	if fee.ID == uuid.Nil {
		return eris.New("")
	}

	if fee.GroupExpenseID == uuid.Nil {
		return eris.New("")
	}

	if fee.Amount.IsZero() {
		return eris.New("amount cannot be zero")
	}

	if len(groupExpense.Participants) < 1 {
		return eris.New("must have participants")
	}

	return nil
}

func (fc *itemizedSplitFeeCalculator) Split(fee entity.OtherFee, groupExpense entity.GroupExpense) []entity.FeeParticipant {
	participantsCount := len(groupExpense.Participants)
	feeParticipants := make([]entity.FeeParticipant, participantsCount)
	rate := fee.Amount.Div(groupExpense.Subtotal)

	for i, expenseParticipant := range groupExpense.Participants {
		feeParticipants[i] = entity.FeeParticipant{
			OtherFeeID:  fee.ID,
			ProfileID:   expenseParticipant.ParticipantProfileID,
			ShareAmount: expenseParticipant.ShareAmount.Mul(rate),
		}
	}

	return feeParticipants
}

func (fc *itemizedSplitFeeCalculator) GetInfo() dto.FeeCalculationMethodInfo {
	return dto.FeeCalculationMethodInfo{
		Name:        fc.method,
		Display:     "Itemized split",
		Description: "Split the fee by a fixed rate applied to each expense items",
	}
}
