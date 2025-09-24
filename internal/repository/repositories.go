package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/go-crud"
)

type GroupExpenseRepository interface {
	crud.Repository[entity.GroupExpense]
	SyncParticipants(ctx context.Context, groupExpenseID uuid.UUID, participants []entity.ExpenseParticipant) error
}

type ExpenseItemRepository interface {
	crud.Repository[entity.ExpenseItem]
	SyncParticipants(ctx context.Context, expenseItemID uuid.UUID, participants []entity.ItemParticipant) error
}

type ExpenseParticipantRepository interface {
	crud.Repository[entity.ExpenseParticipant]
}

type OtherFeeRepository interface {
	crud.Repository[entity.OtherFee]
	SyncParticipants(ctx context.Context, feeID uuid.UUID, participants []entity.FeeParticipant) error
}

type ExpenseBillRepository interface {
	crud.Repository[entity.ExpenseBill]
}

type TaskQueue interface {
	Enqueue(ctx context.Context, task entity.Task) error
	Ping() error
}
