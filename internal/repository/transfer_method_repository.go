package repository

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type transferMethodRepositoryGorm struct {
	db *gorm.DB
}

func NewTransferMethodRepository(db *gorm.DB) TransferMethodRepository {
	return &transferMethodRepositoryGorm{db}
}

func (tmr *transferMethodRepositoryGorm) FindAll(ctx context.Context, spec entity.TransferMethod) ([]entity.TransferMethod, error) {
	var transferMethods []entity.TransferMethod

	db, err := tmr.getGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	if err = db.Scopes(ezutil.WhereBySpec(spec)).Find(&transferMethods).Error; err != nil {
		return nil, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return transferMethods, nil
}

func (tmr *transferMethodRepositoryGorm) FindFirst(ctx context.Context, spec entity.TransferMethod) (entity.TransferMethod, error) {
	var transferMethod entity.TransferMethod

	db, err := tmr.getGormInstance(ctx)
	if err != nil {
		return entity.TransferMethod{}, err
	}

	if err = db.Scopes(ezutil.WhereBySpec(spec)).First(&transferMethod).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.TransferMethod{}, nil
		}
		return entity.TransferMethod{}, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return transferMethod, nil
}

func (tmr *transferMethodRepositoryGorm) getGormInstance(ctx context.Context) (*gorm.DB, error) {
	tx, err := ezutil.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx, nil
	}

	return tmr.db.WithContext(ctx), nil
}
