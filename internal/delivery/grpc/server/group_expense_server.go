package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr-protos/gen/go/groupexpense/v1"
	"github.com/itsLeonB/billsplittr/internal/delivery/grpc/mapper"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
)

type GroupExpenseServer struct {
	groupexpense.UnimplementedGroupExpenseServiceServer
	validate        *validator.Validate
	groupExpenseSvc service.GroupExpenseService
}

func newGroupExpenseServer(
	validate *validator.Validate,
	groupExpenseSvc service.GroupExpenseService,
) groupexpense.GroupExpenseServiceServer {
	return &GroupExpenseServer{
		validate:        validate,
		groupExpenseSvc: groupExpenseSvc,
	}
}

func (ges *GroupExpenseServer) CreateDraft(ctx context.Context, req *groupexpense.CreateDraftRequest) (*groupexpense.CreateDraftResponse, error) {
	creatorProfileID, err := ezutil.Parse[uuid.UUID](req.GetCreatorProfileId())
	if err != nil {
		return nil, err
	}

	payerProfileID, err := ezutil.Parse[uuid.UUID](req.GetPayerProfileId())
	if err != nil {
		return nil, err
	}

	fees := make([]dto.OtherFeeData, 0)
	if req.GetOtherFees() != nil {
		fees, err = ezutil.MapSliceWithError(req.GetOtherFees(), mapper.FromOtherFeeProto)
		if err != nil {
			return nil, err
		}
	}

	request := dto.NewGroupExpenseRequest{
		CreatorProfileID: creatorProfileID,
		PayerProfileID:   payerProfileID,
		TotalAmount:      ezutil.MoneyToDecimal(req.GetTotalAmount()),
		Subtotal:         ezutil.MoneyToDecimal(req.GetSubtotal()),
		Description:      req.GetDescription(),
		Items:            ezutil.MapSlice(req.GetItems(), mapper.FromExpenseItemProto),
		OtherFees:        fees,
	}

	if err = ges.validate.Struct(request); err != nil {
		return nil, err
	}

	groupExpense, err := ges.groupExpenseSvc.CreateDraft(ctx, request)
	if err != nil {
		return nil, err
	}

	response, err := mapper.ToGroupExpenseResponseProto(groupExpense)
	if err != nil {
		return nil, err
	}

	return &groupexpense.CreateDraftResponse{
		GroupExpense: response,
	}, nil
}

func (ges *GroupExpenseServer) GetAllCreated(ctx context.Context, req *groupexpense.GetAllCreatedRequest) (*groupexpense.GetAllCreatedResponse, error) {
	profileID, err := ezutil.Parse[uuid.UUID](req.GetProfileId())
	if err != nil {
		return nil, err
	}

	groupExpenses, err := ges.groupExpenseSvc.GetAllCreated(ctx, profileID)
	if err != nil {
		return nil, err
	}

	responses, err := ezutil.MapSliceWithError(groupExpenses, mapper.ToGroupExpenseResponseProto)
	if err != nil {
		return nil, err
	}

	return &groupexpense.GetAllCreatedResponse{
		GroupExpenses: responses,
	}, nil
}

func (ges *GroupExpenseServer) GetDetails(ctx context.Context, req *groupexpense.GetDetailsRequest) (*groupexpense.GetDetailsResponse, error) {
	id, err := ezutil.Parse[uuid.UUID](req.GetId())
	if err != nil {
		return nil, err
	}

	groupExpense, err := ges.groupExpenseSvc.GetDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	response, err := mapper.ToGroupExpenseResponseProto(groupExpense)
	if err != nil {
		return nil, err
	}

	return &groupexpense.GetDetailsResponse{
		GroupExpense: response,
	}, nil
}

func (ges *GroupExpenseServer) ConfirmDraft(ctx context.Context, req *groupexpense.ConfirmDraftRequest) (*groupexpense.ConfirmDraftResponse, error) {
	id, err := ezutil.Parse[uuid.UUID](req.GetId())
	if err != nil {
		return nil, err
	}

	profileID, err := ezutil.Parse[uuid.UUID](req.GetProfileId())
	if err != nil {
		return nil, err
	}

	groupExpense, err := ges.groupExpenseSvc.ConfirmDraft(ctx, id, profileID)
	if err != nil {
		return nil, err
	}

	response, err := mapper.ToGroupExpenseResponseProto(groupExpense)
	if err != nil {
		return nil, err
	}

	return &groupexpense.ConfirmDraftResponse{
		GroupExpense: response,
	}, nil
}
