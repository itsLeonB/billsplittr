package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr-protos/gen/go/expensebill/v1"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/delivery/grpc/mapper"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/ungerr"
	"google.golang.org/protobuf/types/known/emptypb"
)

type expenseBillServer struct {
	expensebill.UnimplementedExpenseBillServiceServer
	validate       *validator.Validate
	expenseBillSvc service.ExpenseBillService
}

func newExpenseBillServer(
	validate *validator.Validate,
	expenseBillSvc service.ExpenseBillService,
) expensebill.ExpenseBillServiceServer {
	return &expenseBillServer{
		validate:       validate,
		expenseBillSvc: expenseBillSvc,
	}
}

func (ebs *expenseBillServer) Save(ctx context.Context, req *expensebill.SaveRequest) (*expensebill.SaveResponse, error) {
	if req == nil {
		return nil, ungerr.BadRequestError(appconstant.ErrNilRequest)
	}
	bill := req.GetExpenseBill()
	if bill == nil {
		return nil, ungerr.BadRequestError("expense bill is nil")
	}

	creatorProfileID, err := ezutil.Parse[uuid.UUID](bill.GetCreatorProfileId())
	if err != nil {
		return nil, err
	}

	payerProfileID := uuid.Nil
	if bill.GetPayerProfileId() != "" {
		payerProfileID, err = ezutil.Parse[uuid.UUID](bill.GetPayerProfileId())
		if err != nil {
			return nil, err
		}
	}

	expenseID := uuid.Nil
	if bill.GetGroupExpenseId() != "" {
		expenseID, err = ezutil.Parse[uuid.UUID](bill.GetGroupExpenseId())
		if err != nil {
			return nil, err
		}
	}

	request := dto.NewExpenseBillRequest{
		CreatorProfileID: creatorProfileID,
		PayerProfileID:   payerProfileID,
		GroupExpenseID:   expenseID,
		Filename:         bill.GetObjectKey(),
	}

	if err = ebs.validate.Struct(request); err != nil {
		return nil, err
	}

	response, err := ebs.expenseBillSvc.Save(ctx, request)
	if err != nil {
		return nil, err
	}

	mappedResp, err := mapper.ToExpenseBillResponseProto(response)
	if err != nil {
		return nil, err
	}

	return &expensebill.SaveResponse{ExpenseBill: mappedResp}, nil
}

func (ebs *expenseBillServer) GetAllCreated(ctx context.Context, req *expensebill.GetAllCreatedRequest) (*expensebill.GetAllCreatedResponse, error) {
	if req == nil {
		return nil, ungerr.BadRequestError(appconstant.ErrNilRequest)
	}

	creatorProfileID, err := ezutil.Parse[uuid.UUID](req.GetCreatorProfileId())
	if err != nil {
		return nil, err
	}

	responses, err := ebs.expenseBillSvc.GetAllCreated(ctx, creatorProfileID)
	if err != nil {
		return nil, err
	}

	mappedResps, err := ezutil.MapSliceWithError(responses, mapper.ToExpenseBillResponseProto)
	if err != nil {
		return nil, err
	}

	return &expensebill.GetAllCreatedResponse{ExpenseBills: mappedResps}, nil
}

func (ebs *expenseBillServer) Get(ctx context.Context, req *expensebill.GetRequest) (*expensebill.GetResponse, error) {
	if req == nil {
		return nil, ungerr.BadRequestError(appconstant.ErrNilRequest)
	}

	id, err := ezutil.Parse[uuid.UUID](req.GetId())
	if err != nil {
		return nil, err
	}

	response, err := ebs.expenseBillSvc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedResp, err := mapper.ToExpenseBillResponseProto(response)
	if err != nil {
		return nil, err
	}

	return &expensebill.GetResponse{ExpenseBill: mappedResp}, nil
}

func (ebs *expenseBillServer) Delete(ctx context.Context, req *expensebill.DeleteRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, ungerr.BadRequestError(appconstant.ErrNilRequest)
	}

	id, err := ezutil.Parse[uuid.UUID](req.GetId())
	if err != nil {
		return nil, err
	}

	profileID, err := ezutil.Parse[uuid.UUID](req.GetProfileId())
	if err != nil {
		return nil, err
	}

	return nil, ebs.expenseBillSvc.Delete(ctx, id, profileID)
}
