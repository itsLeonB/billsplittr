package entity_test

import (
	"testing"

	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestExpenseBillTableName(t *testing.T) {
	eb := entity.ExpenseBill{}
	assert.Equal(t, "group_expense_bills", eb.TableName())
}
