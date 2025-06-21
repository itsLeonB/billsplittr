package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type debtTransactionRepositoryGorm struct {
	db *gorm.DB
}

func NewDebtTransactionRepository(db *gorm.DB) DebtTransactionRepository {
	return &debtTransactionRepositoryGorm{db}
}

func (dtr *debtTransactionRepositoryGorm) Insert(ctx context.Context, debtTransaction entity.DebtTransaction) (entity.DebtTransaction, error) {
	db, err := dtr.getGormInstance(ctx)
	if err != nil {
		return entity.DebtTransaction{}, err
	}

	if err = db.Create(&debtTransaction).Error; err != nil {
		return entity.DebtTransaction{}, eris.Wrap(err, appconstant.ErrDataInsert)
	}

	return debtTransaction, nil
}

func (dtr *debtTransactionRepositoryGorm) FindAllByProfileID(
	ctx context.Context,
	userProfileID, friendProfileID uuid.UUID,
) ([]entity.DebtTransaction, error) {
	var transactions []entity.DebtTransaction

	db, err := dtr.getGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	err = db.
		Where("lender_profile_id = ? AND borrower_profile_id = ?", userProfileID, friendProfileID).
		Or("lender_profile_id = ? AND borrower_profile_id = ?", friendProfileID, userProfileID).
		Find(&transactions).
		Error

	if err != nil {
		return nil, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return transactions, nil
}

func (dtr *debtTransactionRepositoryGorm) FindAllByUserProfileID(ctx context.Context, userProfileID uuid.UUID) ([]entity.DebtTransaction, error) {
	var transactions []entity.DebtTransaction

	db, err := dtr.getGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	err = db.
		Where("lender_profile_id = ?", userProfileID).
		Or("borrower_profile_id = ?", userProfileID).
		Preload("TransferMethod").
		Find(&transactions).
		Error

	if err != nil {
		return nil, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return transactions, nil
}

func (dtr *debtTransactionRepositoryGorm) getGormInstance(ctx context.Context) (*gorm.DB, error) {
	tx, err := ezutil.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx, nil
	}

	return dtr.db.WithContext(ctx), nil
}
