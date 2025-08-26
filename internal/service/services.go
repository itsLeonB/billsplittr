package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
)

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) error
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error)
	VerifyToken(ctx context.Context, token string) (bool, map[string]any, error)
}

type ProfileService interface {
	GetByID(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error)
	GetNames(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]string, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]dto.ProfileResponse, error)
}

type FriendshipService interface {
	CreateAnonymous(ctx context.Context, request dto.NewAnonymousFriendshipRequest) (dto.FriendshipResponse, error)
	GetAll(ctx context.Context, profileID uuid.UUID) ([]dto.FriendshipResponse, error)
	GetDetails(ctx context.Context, profileID, friendshipID uuid.UUID) (dto.FriendDetailsResponse, error)
	IsFriends(ctx context.Context, profileID1, profileID2 uuid.UUID) (bool, bool, error)
}

type DebtService interface {
	RecordNewTransaction(ctx context.Context, request dto.NewDebtTransactionRequest) (dto.DebtTransactionResponse, error)
	GetTransactions(ctx context.Context, userProfileID uuid.UUID) ([]dto.DebtTransactionResponse, error)
	ProcessConfirmedGroupExpense(ctx context.Context, groupExpenseID uuid.UUID) error
}

type TransferMethodService interface {
	GetAll(ctx context.Context) ([]dto.TransferMethodResponse, error)
}

type GroupExpenseService interface {
	CreateDraft(ctx context.Context, request dto.NewGroupExpenseRequest) (dto.GroupExpenseResponse, error)
	GetAllCreated(ctx context.Context, profileID uuid.UUID) ([]dto.GroupExpenseResponse, error)
	GetDetails(ctx context.Context, id, profileID uuid.UUID) (dto.GroupExpenseResponse, error)
	GetItemDetails(ctx context.Context, groupExpenseID, expenseItemID, profileID uuid.UUID) (dto.ExpenseItemResponse, error)
	UpdateItem(ctx context.Context, profileID uuid.UUID, request dto.UpdateExpenseItemRequest) (dto.ExpenseItemResponse, error)
	ConfirmDraft(ctx context.Context, id, profileID uuid.UUID) (dto.GroupExpenseResponse, error)
	GetFeeCalculationMethods() []dto.FeeCalculationMethodInfo
	UpdateFee(ctx context.Context, profileID uuid.UUID, request dto.UpdateOtherFeeRequest) (dto.OtherFeeResponse, error)
	AddItem(ctx context.Context, profileID uuid.UUID, request dto.NewExpenseItemRequest) (dto.ExpenseItemResponse, error)
	AddFee(ctx context.Context, profileID uuid.UUID, request dto.NewOtherFeeRequest) (dto.OtherFeeResponse, error)
	RemoveItem(ctx context.Context, request dto.DeleteExpenseItemRequest) error
	RemoveFee(ctx context.Context, request dto.DeleteOtherFeeRequest) error
}

type ExpenseBillService interface {
	Upload(ctx context.Context, request dto.NewExpenseBillRequest) error
}
