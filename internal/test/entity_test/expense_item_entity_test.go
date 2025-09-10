package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestExpenseItemTableName(t *testing.T) {
	ei := entity.ExpenseItem{}
	assert.Equal(t, "group_expense_items", ei.TableName())
}

func TestExpenseItemSimpleName(t *testing.T) {
	ei := entity.ExpenseItem{}
	assert.Equal(t, "expense item", ei.SimpleName())
}

func TestExpenseItemTotalAmount(t *testing.T) {
	ei := entity.ExpenseItem{
		Amount:   decimal.NewFromFloat(10.50),
		Quantity: 3,
	}
	expected := decimal.NewFromFloat(31.50)
	assert.True(t, expected.Equal(ei.TotalAmount()))
}

func TestExpenseItemProfileIDs(t *testing.T) {
	profileID1 := uuid.New()
	profileID2 := uuid.New()

	ei := entity.ExpenseItem{
		Participants: []entity.ItemParticipant{
			{ProfileID: profileID1},
			{ProfileID: profileID2},
		},
	}

	profileIDs := ei.ProfileIDs()
	assert.Len(t, profileIDs, 2)
	assert.Contains(t, profileIDs, profileID1)
	assert.Contains(t, profileIDs, profileID2)
}

func TestItemParticipantTableName(t *testing.T) {
	ip := entity.ItemParticipant{}
	assert.Equal(t, "group_expense_item_participants", ip.TableName())
}
