package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
)

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) error
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error)
}

type UserService interface {
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.UserResponse, error)
	GetProfile(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error)
	GetEntityByID(ctx context.Context, id uuid.UUID) (entity.User, error)
}

type FriendshipService interface {
	CreateAnonymous(ctx context.Context, request dto.NewAnonymousFriendshipRequest) (dto.FriendshipResponse, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]dto.FriendshipResponse, error)
	GetDetails(ctx context.Context, userID, friendshipID uuid.UUID) (dto.FriendDetailsResponse, error)
	IsFriends(ctx context.Context, profileID1, profileID2 uuid.UUID) (bool, error)
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
	GetAllCreated(ctx context.Context, userID uuid.UUID) ([]dto.GroupExpenseResponse, error)
	GetDetails(ctx context.Context, id uuid.UUID) (dto.GroupExpenseResponse, error)
	GetItemDetails(ctx context.Context, groupExpenseID, expenseItemID uuid.UUID) (dto.ExpenseItemResponse, error)
	UpdateItem(ctx context.Context, request dto.UpdateExpenseItemRequest) (dto.ExpenseItemResponse, error)
	ConfirmDraft(ctx context.Context, id uuid.UUID) (dto.GroupExpenseResponse, error)
	GetFeeCalculationMethods() []dto.FeeCalculationMethodInfo
	UpdateFee(ctx context.Context, request dto.UpdateOtherFeeRequest) (dto.OtherFeeResponse, error)
	AddItem(ctx context.Context, request dto.NewExpenseItemRequest) (dto.ExpenseItemResponse, error)
	AddFee(ctx context.Context, request dto.NewOtherFeeRequest) (dto.OtherFeeResponse, error)
	RemoveItem(ctx context.Context, request dto.DeleteExpenseItemRequest) error
	RemoveFee(ctx context.Context, request dto.DeleteOtherFeeRequest) error
}

type ExpenseBillService interface {
	Upload(ctx context.Context, request dto.NewExpenseBillRequest) error
}
