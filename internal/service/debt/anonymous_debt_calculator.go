package debt

import (
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
)

type AnonymousDebtCalculator interface {
	GetAction() appconstant.Action
	MapRequestToEntity(request dto.NewDebtTransactionRequest) entity.DebtTransaction
	MapEntityToResponse(debtTransaction entity.DebtTransaction) dto.DebtTransactionResponse
}

var initFuncs = []func() AnonymousDebtCalculator{
	newBorrowingAnonDebtCalculator,
	newLendingAnonDebtCalculator,
}

func NewAnonymousDebtCalculatorStrategies() map[appconstant.Action]AnonymousDebtCalculator {
	strategyMap := make(map[appconstant.Action]AnonymousDebtCalculator)

	for _, initFunc := range initFuncs {
		calculator := initFunc()
		strategyMap[calculator.GetAction()] = calculator
	}

	return strategyMap
}
