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
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

type debtServiceImpl struct {
	transactor                      ezutil.Transactor
	friendshipRepository            repository.FriendshipRepository
	userService                     UserService
	anonymousDebtCalculatorStrategy map[appconstant.Action]debt.AnonymousDebtCalculator
	debtTransactionRepository       repository.DebtTransactionRepository
	transferMethodRepository        repository.TransferMethodRepository
}

func NewDebtService(
	transactor ezutil.Transactor,
	friendshipRepository repository.FriendshipRepository,
	userService UserService,
	debtTransactionRepository repository.DebtTransactionRepository,
	transferMethodRepository repository.TransferMethodRepository,
) DebtService {
	return &debtServiceImpl{
		transactor,
		friendshipRepository,
		userService,
		debt.NewAnonymousDebtCalculatorStrategies(),
		debtTransactionRepository,
		transferMethodRepository,
	}
}

func (ds *debtServiceImpl) RecordNewTransaction(
	ctx context.Context,
	request dto.NewDebtTransactionRequest,
) (dto.DebtTransactionResponse, error) {
	var response dto.DebtTransactionResponse

	if request.Amount.Compare(decimal.Zero) < 1 {
		return dto.DebtTransactionResponse{}, ezutil.ValidationError("amount must be greater than 0")
	}

	transferMethod, err := ds.getTransferMethod(ctx, request.TransferMethodID)
	if err != nil {
		return dto.DebtTransactionResponse{}, err
	}

	err = ds.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		debtTransactions, err := ds.getDebtTransactionsAndPatchRequest(ctx, &request)
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

func (ds *debtServiceImpl) GetTransactions(ctx context.Context, userID uuid.UUID) ([]dto.DebtTransactionResponse, error) {
	user, err := ds.userService.GetEntityByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	transactions, err := ds.debtTransactionRepository.FindAllByUserProfileID(ctx, user.Profile.ID)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(transactions, mapper.GetDebtTransactionSimpleMapper(user.Profile.ID)), nil
}

func (ds *debtServiceImpl) getDebtTransactionsAndPatchRequest(ctx context.Context, request *dto.NewDebtTransactionRequest) ([]entity.DebtTransaction, error) {
	user, err := ds.userService.GetEntityByID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	request.UserProfileID = user.Profile.ID

	friendship, err := ds.getFriendship(ctx, request.UserProfileID, request.FriendProfileID)
	if err != nil {
		return nil, err
	}

	if friendship.Type != appconstant.Anonymous {
		return nil, ezutil.UnprocessableEntityError("flow is forbidden for non-anonymous friendships")
	}

	return ds.debtTransactionRepository.FindAllByProfileID(ctx, request.UserProfileID, request.FriendProfileID)
}

func (ds *debtServiceImpl) getFriendship(ctx context.Context, userProfileID, friendProfileID uuid.UUID) (entity.Friendship, error) {
	friendshipSpec := entity.FriendshipSpecification{}
	userProfile := entity.UserProfile{}
	userProfile.ID = userProfileID
	friendProfile := entity.UserProfile{}
	friendProfile.ID = friendProfileID

	friendshipEntity, err := mapper.OrderProfilesToFriendship(userProfile, friendProfile)
	if err != nil {
		return entity.Friendship{}, err
	}

	friendshipSpec.Friendship = friendshipEntity

	friendship, err := ds.friendshipRepository.FindFirst(ctx, friendshipSpec)
	if err != nil {
		return entity.Friendship{}, err
	}
	if friendship.IsZero() {
		return entity.Friendship{}, ezutil.NotFoundError(appconstant.ErrFriendshipNotFound)
	}
	if friendship.IsDeleted() {
		return entity.Friendship{}, ezutil.UnprocessableEntityError(appconstant.ErrFriendshipDeleted)
	}

	return friendship, nil
}

func (ds *debtServiceImpl) selectAnonCalculator(action appconstant.Action) (debt.AnonymousDebtCalculator, error) {
	calculator, ok := ds.anonymousDebtCalculatorStrategy[action]
	if !ok {
		return nil, eris.Errorf("unsupported anonymous debt calculator action: %s", action)
	}

	return calculator, nil
}

func (ds *debtServiceImpl) getTransferMethod(ctx context.Context, id uuid.UUID) (entity.TransferMethod, error) {
	transferMethodSpec := entity.TransferMethod{}
	transferMethodSpec.ID = id

	transferMethod, err := ds.transferMethodRepository.FindFirst(ctx, transferMethodSpec)
	if err != nil {
		return entity.TransferMethod{}, err
	}

	if transferMethod.IsZero() {
		return entity.TransferMethod{}, ezutil.NotFoundError(fmt.Sprintf(appconstant.ErrTransferMethodNotFound, id))
	}

	if transferMethod.IsDeleted() {
		return entity.TransferMethod{}, ezutil.UnprocessableEntityError(fmt.Sprintf(appconstant.ErrTransferMethodDeleted, id))
	}

	return transferMethod, nil
}
