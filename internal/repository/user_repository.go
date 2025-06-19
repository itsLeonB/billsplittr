package repository

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type userRepositoryGorm struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryGorm{
		db: db,
	}
}

func (ur *userRepositoryGorm) Insert(ctx context.Context, user entity.User) (entity.User, error) {
	db, err := ur.getGormInstance(ctx)
	if err != nil {
		return entity.User{}, err
	}

	err = db.Create(&user).Error
	if err != nil {
		return entity.User{}, eris.Wrap(err, appconstant.ErrDataInsert)
	}

	return user, nil
}

func (ur *userRepositoryGorm) Find(ctx context.Context, spec entity.User) (entity.User, error) {
	var user entity.User

	db, err := ur.getGormInstance(ctx)
	if err != nil {
		return entity.User{}, err
	}

	db = db.Scopes(ezutil.WhereBySpec(spec), ezutil.PreloadRelations([]string{"Profile"}))
	if err := db.First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.User{}, nil // No user found
		}
		return entity.User{}, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return user, nil
}

func (ur *userRepositoryGorm) getGormInstance(ctx context.Context) (*gorm.DB, error) {
	tx, err := ezutil.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx, nil
	}

	return ur.db.WithContext(ctx), nil
}
