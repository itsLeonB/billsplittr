package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/billsplittr/internal/service/debt"
	"github.com/itsLeonB/ezutil/v2"
	crud "github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

type debtServiceImpl struct {
	transactor                      crud.Transactor
	anonymousDebtCalculatorStrategy map[appconstant.Action]debt.AnonymousDebtCalculator
	debtTransactionRepository       repository.DebtTransactionRepository
	transferMethodRepository        repository.TransferMethodRepository
	groupExpenseRepository          repository.GroupExpenseRepository
	friendshipService               FriendshipService
}

func NewDebtService(
	transactor crud.Transactor,
	debtTransactionRepository repository.DebtTransactionRepository,
	transferMethodRepository repository.TransferMethodRepository,
	groupExpenseRepository repository.GroupExpenseRepository,
	friendshipService FriendshipService,
) DebtService {
	return &debtServiceImpl{
		transactor,
		debt.NewAnonymousDebtCalculatorStrategies(),
		debtTransactionRepository,
		transferMethodRepository,
		groupExpenseRepository,
		friendshipService,
	}
}

func (ds *debtServiceImpl) RecordNewTransaction(
	ctx context.Context,
	request dto.NewDebtTransactionRequest,
) (dto.DebtTransactionResponse, error) {
	var response dto.DebtTransactionResponse

	transferMethod, err := ds.getTransferMethod(ctx, request.TransferMethodID)
	if err != nil {
		return dto.DebtTransactionResponse{}, err
	}

	if err = ds.Validate(ctx, request); err != nil {
		return dto.DebtTransactionResponse{}, err
	}

	err = ds.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		debtTransactions, err := ds.debtTransactionRepository.FindAllByProfileID(ctx, request.UserProfileID, request.FriendProfileID)
		if err != nil {
			return err
		}

		calculator, err := ds.selectAnonCalculator(request.Action)
		if err != nil {
			return err
		}

		newDebt := calculator.MapRequestToEntity(request)

		if err = calculator.Validate(newDebt, debtTransactions); err != nil {
			return err
		}

		insertedDebt, err := ds.debtTransactionRepository.Insert(ctx, newDebt)
		if err != nil {
			return err
		}

		insertedDebt.TransferMethod = transferMethod

		response = calculator.MapEntityToResponse(insertedDebt)

		return nil
	})

	if err != nil {
		return dto.DebtTransactionResponse{}, err
	}

	return response, nil
}

func (ds *debtServiceImpl) GetTransactions(ctx context.Context, profileID uuid.UUID) ([]dto.DebtTransactionResponse, error) {
	transactions, err := ds.debtTransactionRepository.FindAllByUserProfileID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(transactions, mapper.GetDebtTransactionSimpleMapper(profileID)), nil
}

func (ds *debtServiceImpl) ProcessConfirmedGroupExpense(ctx context.Context, groupExpenseID uuid.UUID) error {
	return ds.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[entity.GroupExpense]{}
		spec.Model.ID = groupExpenseID
		spec.PreloadRelations = []string{"Participants"}
		spec.ForUpdate = true

		groupExpense, err := ds.groupExpenseRepository.FindFirst(ctx, spec)
		if err != nil {
			return err
		}

		if !groupExpense.ParticipantsConfirmed {
			return ungerr.UnprocessableEntityError(fmt.Sprintf("group expense %s is not confirmed by participants", groupExpenseID))
		}

		transferMethod, err := ds.getGroupExpenseTransferMethod(ctx)
		if err != nil {
			return err
		}

		debtTransactions := mapper.GroupExpenseToDebtTransactions(groupExpense, transferMethod.ID)

		_, err = ds.debtTransactionRepository.BatchInsert(ctx, debtTransactions)

		return err
	})
}

func (ds *debtServiceImpl) Validate(ctx context.Context, req dto.NewDebtTransactionRequest) error {
	if req.Amount.Compare(decimal.Zero) < 1 {
		return ungerr.ValidationError("amount must be greater than 0")
	}

	isFriends, isAnonymous, err := ds.friendshipService.IsFriends(ctx, req.UserProfileID, req.FriendProfileID)
	if err != nil {
		return err
	}
	if !isFriends {
		return ungerr.UnprocessableEntityError("both profiles are not friends")
	}
	if !isAnonymous {
		return ungerr.UnprocessableEntityError("flow is forbidden for non-anonymous friendships")
	}

	return nil
}

func (ds *debtServiceImpl) selectAnonCalculator(action appconstant.Action) (debt.AnonymousDebtCalculator, error) {
	calculator, ok := ds.anonymousDebtCalculatorStrategy[action]
	if !ok {
		return nil, eris.Errorf("unsupported anonymous debt calculator action: %s", action)
	}

	return calculator, nil
}

func (ds *debtServiceImpl) getTransferMethod(ctx context.Context, id uuid.UUID) (entity.TransferMethod, error) {
	spec := crud.Specification[entity.TransferMethod]{}
	spec.Model.ID = id

	transferMethod, err := ds.transferMethodRepository.FindFirst(ctx, spec)
	if err != nil {
		return entity.TransferMethod{}, err
	}

	if transferMethod.IsZero() {
		return entity.TransferMethod{}, ungerr.NotFoundError(fmt.Sprintf(appconstant.ErrTransferMethodNotFound, id))
	}

	if transferMethod.IsDeleted() {
		return entity.TransferMethod{}, ungerr.UnprocessableEntityError(fmt.Sprintf(appconstant.ErrTransferMethodDeleted, id))
	}

	return transferMethod, nil
}

func (ds *debtServiceImpl) getGroupExpenseTransferMethod(ctx context.Context) (entity.TransferMethod, error) {
	spec := crud.Specification[entity.TransferMethod]{}
	spec.Model.Name = appconstant.GroupExpenseTransferMethod

	transferMethod, err := ds.transferMethodRepository.FindFirst(ctx, spec)
	if err != nil {
		return entity.TransferMethod{}, err
	}

	if transferMethod.IsZero() {
		return entity.TransferMethod{}, eris.New("group expense transfer method not found")
	}

	if transferMethod.IsDeleted() {
		return entity.TransferMethod{}, eris.New("group expense transfer method is deleted")
	}

	return transferMethod, nil
}
