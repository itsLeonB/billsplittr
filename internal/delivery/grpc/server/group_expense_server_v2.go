package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr-protos/gen/go/groupexpense/v2"
	"github.com/itsLeonB/billsplittr/internal/delivery/grpc/mapper"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
)

type groupExpenseServerV2 struct {
	groupexpense.UnimplementedGroupExpenseServiceServer
	validate        *validator.Validate
	groupExpenseSvc service.GroupExpenseService
}

func newGroupExpenseServerV2(
	validate *validator.Validate,
	groupExpenseSvc service.GroupExpenseService,
) groupexpense.GroupExpenseServiceServer {
	return &groupExpenseServerV2{
		validate:        validate,
		groupExpenseSvc: groupExpenseSvc,
	}
}

func (ges *groupExpenseServerV2) CreateDraft(ctx context.Context, req *groupexpense.CreateDraftRequest) (*groupexpense.CreateDraftResponse, error) {
	creatorProfileID, err := ezutil.Parse[uuid.UUID](req.GetCreatorProfileId())
	if err != nil {
		return nil, err
	}

	request := dto.NewDraftExpense{
		CreatorProfileID: creatorProfileID,
		Description:      req.GetDescription(),
	}

	if err = ges.validate.Struct(request); err != nil {
		return nil, err
	}

	groupExpense, err := ges.groupExpenseSvc.CreateDraftV2(ctx, request)
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
