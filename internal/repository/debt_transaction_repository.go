package repository

import (
	"context"

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
