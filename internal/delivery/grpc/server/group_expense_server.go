package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr-protos/gen/go/groupexpense/v1"
	"github.com/itsLeonB/billsplittr/internal/delivery/grpc/mapper"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/billsplittr/internal/util/uuidutil"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/rotisserie/eris"
	"google.golang.org/protobuf/types/known/emptypb"
)

type groupExpenseServer struct {
	groupexpense.UnimplementedGroupExpenseServiceServer
	validate        *validator.Validate
	groupExpenseSvc service.GroupExpenseService
}

func newGroupExpenseServer(
	validate *validator.Validate,
	groupExpenseSvc service.GroupExpenseService,
) groupexpense.GroupExpenseServiceServer {
	return &groupExpenseServer{
		validate:        validate,
		groupExpenseSvc: groupExpenseSvc,
	}
}

func (ges *groupExpenseServer) CreateDraft(ctx context.Context, req *groupexpense.CreateDraftRequest) (*groupexpense.CreateDraftResponse, error) {
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

func (ges *groupExpenseServer) GetAllCreated(ctx context.Context, req *groupexpense.GetAllCreatedRequest) (*groupexpense.GetAllCreatedResponse, error) {
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

func (ges *groupExpenseServer) GetDetails(ctx context.Context, req *groupexpense.GetDetailsRequest) (*groupexpense.GetDetailsResponse, error) {
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

func (ges *groupExpenseServer) ConfirmDraft(ctx context.Context, req *groupexpense.ConfirmDraftRequest) (*groupexpense.ConfirmDraftResponse, error) {
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

func (ges *groupExpenseServer) Delete(ctx context.Context, req *groupexpense.DeleteRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, eris.New("request is nil")
	}

	id, err := ezutil.Parse[uuid.UUID](req.GetId())
	if err != nil {
		return nil, err
	}

	profileID, err := ezutil.Parse[uuid.UUID](req.GetProfileId())
	if err != nil {
		return nil, err
	}

	return nil, ges.groupExpenseSvc.Delete(ctx, id, profileID)
}

func (ges *groupExpenseServer) SyncParticipants(ctx context.Context, req *groupexpense.SyncParticipantsRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, eris.New("request is nil")
	}

	participantProfileIDs, err := ezutil.MapSliceWithError(req.GetParticipantProfileIds(), uuidutil.Parse)
	if err != nil {
		return nil, err
	}

	payerProfileID, err := uuidutil.Parse(req.GetPayerProfileId())
	if err != nil {
		return nil, err
	}

	userProfileID, err := uuidutil.Parse(req.GetUserProfileId())
	if err != nil {
		return nil, err
	}

	expenseID, err := uuidutil.Parse(req.GetGroupExpenseId())
	if err != nil {
		return nil, err
	}

	request := dto.ExpenseParticipantsRequest{
		ParticipantProfileIDs: participantProfileIDs,
		PayerProfileID:        payerProfileID,
		UserProfileID:         userProfileID,
		GroupExpenseID:        expenseID,
	}

	return nil, ges.groupExpenseSvc.SyncParticipants(ctx, request)
}
