package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/entity"
	crud "github.com/itsLeonB/go-crud"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type debtTransactionRepositoryGorm struct {
	crud.CRUDRepository[entity.DebtTransaction]
	db *gorm.DB
}

func NewDebtTransactionRepository(db *gorm.DB) DebtTransactionRepository {
	return &debtTransactionRepositoryGorm{
		crud.NewCRUDRepository[entity.DebtTransaction](db),
		db,
	}
}

func (dtr *debtTransactionRepositoryGorm) FindAllByProfileID(
	ctx context.Context,
	userProfileID, friendProfileID uuid.UUID,
) ([]entity.DebtTransaction, error) {
	var transactions []entity.DebtTransaction

	db, err := dtr.GetGormInstance(ctx)
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

	db, err := dtr.GetGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	err = db.
		Where("lender_profile_id = ?", userProfileID).
		Or("borrower_profile_id = ?", userProfileID).
		Preload("TransferMethod").
		Scopes(crud.DefaultOrder()).
		Find(&transactions).
		Error

	if err != nil {
		return nil, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return transactions, nil
}
