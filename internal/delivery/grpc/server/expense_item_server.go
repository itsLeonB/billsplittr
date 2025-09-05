package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr-protos/gen/go/domain/v1"
	"github.com/itsLeonB/billsplittr-protos/gen/go/expenseitem/v1"
	"github.com/itsLeonB/billsplittr/internal/delivery/grpc/mapper"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/ungerr"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ExpenseItemServer struct {
	expenseitem.UnimplementedExpenseItemServiceServer
	validate       *validator.Validate
	expenseItemSvc service.ExpenseItemService
}

func newExpenseItemServer(
	validate *validator.Validate,
	expenseItemSvc service.ExpenseItemService,
) expenseitem.ExpenseItemServiceServer {
	return &ExpenseItemServer{
		validate:       validate,
		expenseItemSvc: expenseItemSvc,
	}
}

func (eis *ExpenseItemServer) Add(ctx context.Context, req *expenseitem.AddRequest) (*expenseitem.AddResponse, error) {
	profileID, err := ezutil.Parse[uuid.UUID](req.GetProfileId())
	if err != nil {
		return nil, err
	}

	groupExpenseID, err := ezutil.Parse[uuid.UUID](req.GetGroupExpenseId())
	if err != nil {
		return nil, err
	}

	if req.GetExpenseItem() == nil {
		return nil, ungerr.BadRequestError("expense item data is empty")
	}

	request := dto.NewExpenseItemRequest{
		ProfileID:       profileID,
		GroupExpenseID:  groupExpenseID,
		ExpenseItemData: mapper.FromExpenseItemProto(req.GetExpenseItem()),
	}

	if err = eis.validate.Struct(request); err != nil {
		return nil, err
	}

	expenseItem, err := eis.expenseItemSvc.Add(ctx, request)
	if err != nil {
		return nil, err
	}

	return &expenseitem.AddResponse{
		ExpenseItem: mapper.ToExpenseItemResponseProto(expenseItem),
	}, nil
}

func (eis *ExpenseItemServer) GetDetails(ctx context.Context, req *expenseitem.GetDetailsRequest) (*expenseitem.GetDetailsResponse, error) {
	id, err := ezutil.Parse[uuid.UUID](req.GetId())
	if err != nil {
		return nil, err
	}

	groupExpenseID, err := ezutil.Parse[uuid.UUID](req.GetGroupExpenseId())
	if err != nil {
		return nil, err
	}

	expenseItem, err := eis.expenseItemSvc.GetDetails(ctx, groupExpenseID, id)
	if err != nil {
		return nil, err
	}

	return &expenseitem.GetDetailsResponse{
		ExpenseItem: mapper.ToExpenseItemResponseProto(expenseItem),
	}, nil
}

func (eis *ExpenseItemServer) Update(ctx context.Context, req *expenseitem.UpdateRequest) (*expenseitem.UpdateResponse, error) {
	profileID, err := ezutil.Parse[uuid.UUID](req.GetProfileId())
	if err != nil {
		return nil, err
	}

	id, err := ezutil.Parse[uuid.UUID](req.GetId())
	if err != nil {
		return nil, err
	}

	groupExpenseID, err := ezutil.Parse[uuid.UUID](req.GetGroupExpenseId())
	if err != nil {
		return nil, err
	}

	if req.GetExpenseItem() == nil {
		return nil, ungerr.BadRequestError("expense item data is nil")
	}

	participantMapFunc := func(participant *domain.ItemParticipant) (dto.ItemParticipantData, error) {
		profileID, err := ezutil.Parse[uuid.UUID](participant.GetProfileId())
		if err != nil {
			return dto.ItemParticipantData{}, err
		}

		return dto.ItemParticipantData{
			ProfileID: profileID,
			Share:     decimal.NewFromFloat32(participant.GetShare()),
		}, nil
	}

	participants, err := ezutil.MapSliceWithError(req.GetExpenseItem().GetParticipants(), participantMapFunc)
	if err != nil {
		return nil, err
	}

	request := dto.UpdateExpenseItemRequest{
		ProfileID:       profileID,
		ID:              id,
		GroupExpenseID:  groupExpenseID,
		ExpenseItemData: mapper.FromExpenseItemProto(req.GetExpenseItem()),
		Participants:    participants,
	}

	if err = eis.validate.Struct(request); err != nil {
		return nil, err
	}

	expenseItem, err := eis.expenseItemSvc.Update(ctx, request)
	if err != nil {
		return nil, err
	}

	return &expenseitem.UpdateResponse{
		ExpenseItem: mapper.ToExpenseItemResponseProto(expenseItem),
	}, nil
}

func (eis *ExpenseItemServer) Remove(ctx context.Context, req *expenseitem.RemoveRequest) (*emptypb.Empty, error) {
	profileID, err := ezutil.Parse[uuid.UUID](req.GetProfileId())
	if err != nil {
		return nil, err
	}

	id, err := ezutil.Parse[uuid.UUID](req.GetId())
	if err != nil {
		return nil, err
	}

	groupExpenseID, err := ezutil.Parse[uuid.UUID](req.GetGroupExpenseId())
	if err != nil {
		return nil, err
	}

	err = eis.expenseItemSvc.Remove(ctx, profileID, id, groupExpenseID)

	return nil, err
}
