package constants_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/stretchr/testify/assert"
)

func TestImageTypes_Contains(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    bool
	}{
		{"JPEG image", "image/jpeg", true},
		{"JPG image", "image/jpg", true},
		{"PNG image", "image/png", true},
		{"Text file", "text/plain", false},
		{"PDF file", "application/pdf", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, exists := appconstant.ImageTypes[tt.contentType]
			assert.Equal(t, tt.expected, exists)
		})
	}
}

func TestDebtTransactionType_Constants(t *testing.T) {
	assert.Equal(t, appconstant.DebtTransactionType("LEND"), appconstant.Lend)
	assert.Equal(t, appconstant.DebtTransactionType("REPAY"), appconstant.Repay)
}

func TestFriendshipType_Constants(t *testing.T) {
	assert.Equal(t, appconstant.FriendshipType("REAL"), appconstant.Real)
	assert.Equal(t, appconstant.FriendshipType("ANON"), appconstant.Anonymous)
}

func TestFeeCalculationMethod_Constants(t *testing.T) {
	assert.Equal(t, appconstant.FeeCalculationMethod("EQUAL_SPLIT"), appconstant.EqualSplitFee)
	assert.Equal(t, appconstant.FeeCalculationMethod("ITEMIZED_SPLIT"), appconstant.ItemizedSplitFee)
}

func TestGroupExpenseTransferMethod_Constant(t *testing.T) {
	assert.Equal(t, "GROUP_EXPENSE", appconstant.GroupExpenseTransferMethod)
}
