package service_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/service/fee"
	"github.com/stretchr/testify/assert"
)

func TestNewFeeCalculatorRegistry(t *testing.T) {
	registry := fee.NewFeeCalculatorRegistry()

	assert.NotNil(t, registry)
	assert.Len(t, registry, 2)
	assert.Contains(t, registry, appconstant.EqualSplitFee)
	assert.Contains(t, registry, appconstant.ItemizedSplitFee)

	equalSplitCalc := registry[appconstant.EqualSplitFee]
	assert.NotNil(t, equalSplitCalc)
	assert.Equal(t, appconstant.EqualSplitFee, equalSplitCalc.GetMethod())

	itemizedSplitCalc := registry[appconstant.ItemizedSplitFee]
	assert.NotNil(t, itemizedSplitCalc)
	assert.Equal(t, appconstant.ItemizedSplitFee, itemizedSplitCalc.GetMethod())
}
