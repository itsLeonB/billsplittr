package repository

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/util"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type groupExpenseRepositoryGorm struct {
	db *gorm.DB
}

func NewGroupExpenseRepository(db *gorm.DB) GroupExpenseRepository {
	return &groupExpenseRepositoryGorm{db}
}

func (ger *groupExpenseRepositoryGorm) Insert(ctx context.Context, groupExpense entity.GroupExpense) (entity.GroupExpense, error) {
	db, err := ger.getGormInstance(ctx)
	if err != nil {
		return entity.GroupExpense{}, err
	}

	if err = db.Create(&groupExpense).Error; err != nil {
		return entity.GroupExpense{}, eris.Wrap(err, appconstant.ErrDataInsert)
	}

	return groupExpense, nil
}

func (ger *groupExpenseRepositoryGorm) FindAll(ctx context.Context, spec entity.GroupExpense) ([]entity.GroupExpense, error) {
	var groupExpenses []entity.GroupExpense

	db, err := ger.getGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	err = db.Scopes(ezutil.WhereBySpec(spec), util.DefaultOrder()).
		Find(&groupExpenses).
		Error

	if err != nil {
		return nil, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return groupExpenses, nil
}

func (ger *groupExpenseRepositoryGorm) FindFirst(ctx context.Context, spec entity.GroupExpenseSpecification) (entity.GroupExpense, error) {
	var groupExpense entity.GroupExpense

	db, err := ger.getGormInstance(ctx)
	if err != nil {
		return entity.GroupExpense{}, err
	}

	err = db.Scopes(ezutil.WhereBySpec(spec.GroupExpense), ezutil.PreloadRelations(spec.PreloadRelations)).
		Take(&groupExpense).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.GroupExpense{}, nil
		}
		return entity.GroupExpense{}, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return groupExpense, nil
}

func (ger *groupExpenseRepositoryGorm) getGormInstance(ctx context.Context) (*gorm.DB, error) {
	tx, err := ezutil.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx, nil
	}
	return ger.db.WithContext(ctx), nil
}
