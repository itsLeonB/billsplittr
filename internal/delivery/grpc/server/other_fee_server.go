package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr-protos/gen/go/otherfee/v1"
	"github.com/itsLeonB/billsplittr/internal/delivery/grpc/mapper"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/ungerr"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OtherFeeServer struct {
	otherfee.UnimplementedOtherFeeServiceServer
	validate    *validator.Validate
	otherFeeSvc service.OtherFeeService
}

func newOtherFeeServer(
	validate *validator.Validate,
	otherFeeSvc service.OtherFeeService,
) otherfee.OtherFeeServiceServer {
	return &OtherFeeServer{
		validate:    validate,
		otherFeeSvc: otherFeeSvc,
	}
}

func (ofs *OtherFeeServer) Add(ctx context.Context, req *otherfee.AddRequest) (*otherfee.AddResponse, error) {
	profileID, err := ezutil.Parse[uuid.UUID](req.GetProfileId())
	if err != nil {
		return nil, err
	}

	groupExpenseID, err := ezutil.Parse[uuid.UUID](req.GetGroupExpenseId())
	if err != nil {
		return nil, err
	}

	if req.GetOtherFee() == nil {
		return nil, ungerr.BadRequestError("other fee data is nil")
	}

	otherFeeData, err := mapper.FromOtherFeeProto(req.GetOtherFee())
	if err != nil {
		return nil, err
	}

	request := dto.NewOtherFeeRequest{
		ProfileID:      profileID,
		GroupExpenseID: groupExpenseID,
		OtherFeeData:   otherFeeData,
	}

	if err = ofs.validate.Struct(request); err != nil {
		return nil, err
	}

	addedFee, err := ofs.otherFeeSvc.Add(ctx, request)
	if err != nil {
		return nil, err
	}

	response, err := mapper.ToOtherFeeResponseProto(addedFee)
	if err != nil {
		return nil, err
	}

	return &otherfee.AddResponse{
		OtherFee: response,
	}, nil
}

func (ofs *OtherFeeServer) Update(ctx context.Context, req *otherfee.UpdateRequest) (*otherfee.UpdateResponse, error) {
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

	if req.GetOtherFee() == nil {
		return nil, ungerr.BadRequestError("other fee data is nil")
	}

	otherFeeData, err := mapper.FromOtherFeeProto(req.GetOtherFee())
	if err != nil {
		return nil, err
	}

	request := dto.UpdateOtherFeeRequest{
		ProfileID:      profileID,
		ID:             id,
		GroupExpenseID: groupExpenseID,
		OtherFeeData:   otherFeeData,
	}

	if err = ofs.validate.Struct(request); err != nil {
		return nil, err
	}

	otherFee, err := ofs.otherFeeSvc.Update(ctx, request)
	if err != nil {
		return nil, err
	}

	response, err := mapper.ToOtherFeeResponseProto(otherFee)
	if err != nil {
		return nil, err
	}

	return &otherfee.UpdateResponse{
		OtherFee: response,
	}, nil
}

func (ofs *OtherFeeServer) Remove(ctx context.Context, req *otherfee.RemoveRequest) (*emptypb.Empty, error) {
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

	err = ofs.otherFeeSvc.Remove(ctx, profileID, id, groupExpenseID)

	return nil, err
}

func (ofs *OtherFeeServer) GetCalculationMethods(ctx context.Context, _ *emptypb.Empty) (*otherfee.GetCalculationMethodsResponse, error) {
	responses, err := ezutil.MapSliceWithError(ofs.otherFeeSvc.GetCalculationMethods(), mapper.ToCalculationMethodInfoProto)
	if err != nil {
		return nil, err
	}

	return &otherfee.GetCalculationMethodsResponse{
		Methods: responses,
	}, nil
}
