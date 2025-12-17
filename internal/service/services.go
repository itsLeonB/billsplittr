package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
)

type GroupExpenseService interface {
	// V1
	CreateDraft(ctx context.Context, request dto.NewGroupExpenseRequest) (dto.GroupExpenseResponse, error)
	GetAllCreated(ctx context.Context, profileID uuid.UUID) ([]dto.GroupExpenseResponse, error)
	GetDetails(ctx context.Context, id uuid.UUID) (dto.GroupExpenseResponse, error)
	ConfirmDraft(ctx context.Context, id, profileID uuid.UUID) (dto.GroupExpenseResponse, error)
	GetUnconfirmedGroupExpenseForUpdate(ctx context.Context, profileID, id uuid.UUID) (entity.GroupExpense, error)
	ParseFromBillText(ctx context.Context) error

	// V2
	CreateDraftV2(ctx context.Context, req dto.NewDraftExpense) (dto.GroupExpenseResponse, error)
}

type ExpenseItemService interface {
	Add(ctx context.Context, request dto.NewExpenseItemRequest) (dto.ExpenseItemResponse, error)
	GetDetails(ctx context.Context, groupExpenseID, expenseItemID uuid.UUID) (dto.ExpenseItemResponse, error)
	Update(ctx context.Context, request dto.UpdateExpenseItemRequest) (dto.ExpenseItemResponse, error)
	Remove(ctx context.Context, profileID, id, groupExpenseID uuid.UUID) error
}

type OtherFeeService interface {
	Add(ctx context.Context, request dto.NewOtherFeeRequest) (dto.OtherFeeResponse, error)
	Update(ctx context.Context, request dto.UpdateOtherFeeRequest) (dto.OtherFeeResponse, error)
	Remove(ctx context.Context, profileID, id, groupExpenseID uuid.UUID) error
	GetCalculationMethods() []dto.FeeCalculationMethodInfo
}

type ExpenseBillService interface {
	Save(ctx context.Context, req dto.NewExpenseBillRequest) (dto.ExpenseBillResponse, error)
	GetAllCreated(ctx context.Context, creatorProfileID uuid.UUID) ([]dto.ExpenseBillResponse, error)
	Get(ctx context.Context, id uuid.UUID) (dto.ExpenseBillResponse, error)
	Delete(ctx context.Context, id, profileID uuid.UUID) error
	EnqueueCleanup(ctx context.Context) error
}

type LLMService interface {
	Prompt(ctx context.Context, systemMsg, userMsg string) (string, error)
}
