package fee

import (
	"log"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
)

var namespace = "[FeeCalculator]"

type FeeCalculator interface {
	GetMethod() appconstant.FeeCalculationMethod
	Validate(fee entity.OtherFee, groupExpense entity.GroupExpense) error
	Split(fee entity.OtherFee, groupExpense entity.GroupExpense) []entity.FeeParticipant
	GetInfo() dto.FeeCalculationMethodInfo
}

var initFuncs = []func() FeeCalculator{
	newEqualSplitFeeCalculator,
	newItemizedSplitFeeCalculator,
}

func NewFeeCalculatorRegistry() map[appconstant.FeeCalculationMethod]FeeCalculator {
	registry := make(map[appconstant.FeeCalculationMethod]FeeCalculator)

	for _, initFunc := range initFuncs {
		if initFunc == nil {
			log.Fatalf("%s initFunc is nil", namespace)
		}

		calculator := initFunc()
		if calculator == nil {
			log.Fatalf("%s calculator is nil", namespace)
		}

		method := calculator.GetMethod()
		if _, exists := registry[method]; exists {
			log.Fatalf("%s duplicate calculator for method: %s", namespace, method)
		}

		registry[calculator.GetMethod()] = calculator
	}

	return registry
}
