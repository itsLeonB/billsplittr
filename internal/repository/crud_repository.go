package repository

import (
	"context"
	"reflect"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/util"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type crudRepositoryGorm[T any] struct {
	db *gorm.DB
}

func NewCRUDRepository[T any](db *gorm.DB) CRUDRepository[T] {
	return &crudRepositoryGorm[T]{db}
}

func (cr *crudRepositoryGorm[T]) Insert(ctx context.Context, model T) (T, error) {
	var zero T

	if err := cr.checkZeroValue(model); err != nil {
		return zero, err
	}

	db, err := cr.GetGormInstance(ctx)
	if err != nil {
		return zero, err
	}

	if err = db.Create(&model).Error; err != nil {
		return zero, eris.Wrap(err, appconstant.ErrDataInsert)
	}

	return model, nil
}

func (cr *crudRepositoryGorm[T]) FindAll(ctx context.Context, spec entity.GenericSpec[T]) ([]T, error) {
	var models []T

	db, err := cr.GetGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	err = db.Scopes(
		ezutil.WhereBySpec(spec.Model),
		util.DefaultOrder(),
		ezutil.PreloadRelations(spec.PreloadRelations),
		util.ForUpdate(spec.ForUpdate),
	).
		Find(&models).
		Error

	if err != nil {
		return nil, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return models, nil
}

func (cr *crudRepositoryGorm[T]) FindFirst(ctx context.Context, spec entity.GenericSpec[T]) (T, error) {
	var model T

	db, err := cr.GetGormInstance(ctx)
	if err != nil {
		return model, err
	}

	err = db.Scopes(
		ezutil.WhereBySpec(spec.Model),
		util.DefaultOrder(),
		ezutil.PreloadRelations(spec.PreloadRelations),
		util.ForUpdate(spec.ForUpdate),
	).
		First(&model).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model, nil
		}
		return model, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return model, nil
}

func (cr *crudRepositoryGorm[T]) Update(ctx context.Context, model T) (T, error) {
	var zero T

	if err := cr.checkZeroValue(model); err != nil {
		return zero, err
	}

	db, err := cr.GetGormInstance(ctx)
	if err != nil {
		return zero, err
	}

	if err = db.Save(&model).Error; err != nil {
		return zero, eris.Wrap(err, appconstant.ErrDataUpdate)
	}

	return model, nil
}

func (cr *crudRepositoryGorm[T]) Delete(ctx context.Context, model T) error {
	if err := cr.checkZeroValue(model); err != nil {
		return err
	}

	db, err := cr.GetGormInstance(ctx)
	if err != nil {
		return err
	}

	if err = db.Unscoped().Delete(&model).Error; err != nil {
		return eris.Wrap(err, appconstant.ErrDataDelete)
	}

	return nil
}

func (cr *crudRepositoryGorm[T]) checkZeroValue(model T) error {
	if reflect.DeepEqual(model, *new(T)) {
		return eris.New("model cannot be zero value")
	}

	return nil
}

func (cr *crudRepositoryGorm[T]) GetGormInstance(ctx context.Context) (*gorm.DB, error) {
	tx, err := ezutil.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx, nil
	}

	return cr.db.WithContext(ctx), nil
}
