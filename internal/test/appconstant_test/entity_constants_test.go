package appconstant_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/stretchr/testify/assert"
)

func TestFeeCalculationMethodConstants(t *testing.T) {
	assert.Equal(t, appconstant.FeeCalculationMethod("EQUAL_SPLIT"), appconstant.EqualSplitFee)
	assert.Equal(t, appconstant.FeeCalculationMethod("ITEMIZED_SPLIT"), appconstant.ItemizedSplitFee)
}

func TestFeeCalculationMethodType(t *testing.T) {
	method := appconstant.EqualSplitFee
	assert.IsType(t, appconstant.FeeCalculationMethod(""), method)
}
