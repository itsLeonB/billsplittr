package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/drex-protos/gen/go/debt/v1"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"
)

type debtServiceImpl struct {
	debtClient        debt.DebtServiceClient
	friendshipService FriendshipService
}

func NewDebtService(
	debtClient debt.DebtServiceClient,
	friendshipService FriendshipService,
) DebtService {
	return &debtServiceImpl{
		debtClient,
		friendshipService,
	}
}

func (ds *debtServiceImpl) RecordNewTransaction(ctx context.Context, req dto.NewDebtTransactionRequest) (dto.DebtTransactionResponse, error) {
	if req.Amount.Compare(decimal.Zero) < 1 {
		return dto.DebtTransactionResponse{}, ungerr.ValidationError("amount must be greater than 0")
	}

	isFriends, isAnonymous, err := ds.friendshipService.IsFriends(ctx, req.UserProfileID, req.FriendProfileID)
	if err != nil {
		return dto.DebtTransactionResponse{}, err
	}
	if !isFriends {
		return dto.DebtTransactionResponse{}, ungerr.UnprocessableEntityError("both profiles are not friends")
	}
	if !isAnonymous {
		return dto.DebtTransactionResponse{}, ungerr.UnprocessableEntityError("flow is forbidden for non-anonymous friendships")
	}

	request := &debt.RecordNewTransactionRequest{
		UserProfileId:    req.UserProfileID.String(),
		FriendProfileId:  req.FriendProfileID.String(),
		Action:           mapper.ToProtoTransactionAction(req.Action),
		Amount:           ezutil.DecimalToMoneyRounded(req.Amount, currency.IDR.String()),
		TransferMethodId: req.TransferMethodID.String(),
		Description:      req.Description,
	}

	response, err := ds.debtClient.RecordNewTransaction(ctx, request)
	if err != nil {
		return dto.DebtTransactionResponse{}, eris.Wrap(err, appconstant.ErrServiceClient)
	}

	return mapper.FromProtoDebtTransactionResponse(response.GetTransaction())
}

func (ds *debtServiceImpl) GetTransactions(ctx context.Context, profileID uuid.UUID) ([]dto.DebtTransactionResponse, error) {
	request := &debt.GetTransactionsRequest{
		UserProfileId: profileID.String(),
	}

	response, err := ds.debtClient.GetTransactions(ctx, request)
	if err != nil {
		return nil, eris.Wrap(err, appconstant.ErrServiceClient)
	}

	return ezutil.MapSliceWithError(response.GetTransactions(), mapper.FromProtoDebtTransactionResponse)
}

func (ds *debtServiceImpl) ProcessConfirmedGroupExpense(ctx context.Context, groupExpense entity.GroupExpense) error {
	if !groupExpense.ParticipantsConfirmed {
		return eris.New("participants are not confirmed")
	}

	mapFunc := func(participant entity.ExpenseParticipant) *debt.ExpenseParticipantData {
		return &debt.ExpenseParticipantData{
			ProfileId:   participant.ParticipantProfileID.String(),
			ShareAmount: ezutil.DecimalToMoneyRounded(participant.ShareAmount, currency.IDR.String()),
		}
	}

	request := &debt.ProcessConfirmedGroupExpenseRequest{
		GroupExpense: &debt.GroupExpenseData{
			Id:               groupExpense.ID.String(),
			PayerProfileId:   groupExpense.PayerProfileID.String(),
			CreatorProfileId: groupExpense.CreatorProfileID.String(),
			Description:      groupExpense.Description,
			Participants:     ezutil.MapSlice(groupExpense.Participants, mapFunc),
		},
	}

	if _, err := ds.debtClient.ProcessConfirmedGroupExpense(ctx, request); err != nil {
		return eris.Wrap(err, appconstant.ErrServiceClient)
	}

	return nil
}

func (ds *debtServiceImpl) GetAllByProfileIDs(ctx context.Context, userProfileID, friendProfileID uuid.UUID) ([]dto.DebtTransactionResponse, error) {
	request := &debt.GetAllByProfileIdsRequest{
		UserProfileId:   userProfileID.String(),
		FriendProfileId: friendProfileID.String(),
	}

	response, err := ds.debtClient.GetAllByProfileIds(ctx, request)
	if err != nil {
		return nil, eris.Wrap(err, appconstant.ErrServiceClient)
	}

	return ezutil.MapSliceWithError(response.GetTransactions(), mapper.FromProtoDebtTransactionResponse)
}
